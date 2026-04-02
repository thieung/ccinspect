package copier

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// skipDirs are skill subdirectories to ignore during listing and copying.
var skipDirs = map[string]bool{
	"_shared": true, ".venv": true, "__pycache__": true,
}

// CopyResult holds information about a single copy operation.
type CopyResult struct {
	Name   string `json:"name"`
	From   string `json:"from"`
	To     string `json:"to"`
	Status string `json:"status"` // "copied", "skipped", "error"
	Detail string `json:"detail,omitempty"`
}

// expandPath resolves "global" to ~/.claude/ and expands ~ in paths.
// Returns (claudePath, isGlobal, error).
func expandPath(path string) (string, bool, error) {
	if path == "global" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", false, fmt.Errorf("cannot determine home directory: %w", err)
		}
		return filepath.Join(home, ".claude"), true, nil
	}

	expanded := expandHome(path)
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", false, fmt.Errorf("invalid path %q: %w", path, err)
	}

	// Check if the path itself is a .claude directory
	if filepath.Base(abs) == ".claude" {
		return abs, false, nil
	}

	// Otherwise assume it's a project root containing .claude/
	claudePath := filepath.Join(abs, ".claude")
	if info, err := os.Stat(claudePath); err == nil && info.IsDir() {
		return claudePath, false, nil
	}

	return "", false, fmt.Errorf("no .claude/ directory found at %s", abs)
}

func expandHome(path string) string {
	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~"+string(filepath.Separator)) {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// subdirForEntity returns the subdirectory name within .claude/ for an entity type.
func subdirForEntity(entityType string) (string, error) {
	switch entityType {
	case "skills":
		return "skills", nil
	case "agents":
		return "agents", nil
	case "commands":
		return "commands", nil
	default:
		return "", fmt.Errorf("unsupported entity type for file-based copy: %s", entityType)
	}
}

// ListAvailable returns the names of available entities of a type in a source path.
func ListAvailable(fromPath string, entityType string) ([]string, error) {
	claudePath, _, err := expandPath(fromPath)
	if err != nil {
		return nil, err
	}

	if entityType == "hooks" {
		return listAvailableHookEvents(claudePath, fromPath)
	}

	subdir, err := subdirForEntity(entityType)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(claudePath, subdir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("no %s found in %s", entityType, fromPath)
	}

	var names []string
	for _, e := range entries {
		if entityType == "skills" {
			// Skills are directories containing SKILL.md
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") && !skipDirs[e.Name()] {
				names = append(names, e.Name())
			}
		} else {
			// Agents, commands are .md files
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				names = append(names, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
	}
	return names, nil
}

func listAvailableHookEvents(claudePath, label string) ([]string, error) {
	settings, err := readSettings(claudePath, label)
	if err != nil {
		return nil, err
	}

	hooksRaw, ok := settings["hooks"]
	if !ok {
		return nil, fmt.Errorf("no hooks found in %s", label)
	}
	hooksMap, ok := hooksRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid hooks format in %s", label)
	}

	var events []string
	for event := range hooksMap {
		events = append(events, event)
	}
	return events, nil
}

// readSettings reads the appropriate settings file for a .claude path.
func readSettings(claudePath, label string) (map[string]any, error) {
	// Try settings.json for global, settings.local.json for projects
	filenames := []string{"settings.json", "settings.local.json"}
	for _, fn := range filenames {
		p := filepath.Join(claudePath, fn)
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			continue
		}
		return m, nil
	}
	return nil, fmt.Errorf("no settings file found in %s", label)
}

// settingsFile returns the correct settings filename for dest.
func settingsFile(isGlobal bool) string {
	if isGlobal {
		return "settings.json"
	}
	return "settings.local.json"
}

// CopyEntities copies entities of entityType from source to destination.
// names can be specific entity names or ["all"] to copy everything.
// force=true overwrites existing entities.
// dryRun=true only shows what would be copied.
func CopyEntities(fromPath, toPath string, entityType string, names []string, force, dryRun bool) ([]CopyResult, error) {
	if entityType == "hooks" {
		return copyHooks(fromPath, toPath, names, force, dryRun)
	}
	return copyFileBased(fromPath, toPath, entityType, names, force, dryRun)
}

func copyFileBased(fromPath, toPath string, entityType string, names []string, force, dryRun bool) ([]CopyResult, error) {
	srcClaude, _, err := expandPath(fromPath)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}

	dstClaude, _, err := expandPath(toPath)
	if err != nil {
		// If the destination doesn't exist yet, create it
		expanded := expandHome(toPath)
		abs, absErr := filepath.Abs(expanded)
		if absErr != nil {
			return nil, fmt.Errorf("destination: %w", err)
		}
		dstClaude = filepath.Join(abs, ".claude")
	}

	subdir, err := subdirForEntity(entityType)
	if err != nil {
		return nil, err
	}

	srcDir := filepath.Join(srcClaude, subdir)
	dstDir := filepath.Join(dstClaude, subdir)

	// Resolve which names to copy
	toCopy, err := resolveNames(srcDir, entityType, names)
	if err != nil {
		return nil, err
	}

	var results []CopyResult

	for _, name := range toCopy {
		result := CopyResult{
			Name: name,
			From: fromPath,
			To:   toPath,
		}

		if entityType == "skills" {
			// Skills are directories
			srcPath := filepath.Join(srcDir, name)
			dstPath := filepath.Join(dstDir, name)

			if !dirExists(srcPath) {
				result.Status = "error"
				result.Detail = fmt.Sprintf("skill %q not found in source", name)
				results = append(results, result)
				continue
			}

			if dirExists(dstPath) && !force {
				result.Status = "skipped"
				result.Detail = "already exists (use --force to overwrite)"
				results = append(results, result)
				continue
			}

			if dryRun {
				result.Status = "would copy"
				result.Detail = fmt.Sprintf("%s → %s", shortenPath(srcPath), shortenPath(dstPath))
				results = append(results, result)
				continue
			}

			if err := copyDir(srcPath, dstPath); err != nil {
				result.Status = "error"
				result.Detail = err.Error()
			} else {
				result.Status = "copied"
				result.Detail = fmt.Sprintf("→ %s", shortenPath(dstPath))
			}
		} else {
			// Agents, commands are .md files
			srcFile := filepath.Join(srcDir, name+".md")
			dstFile := filepath.Join(dstDir, name+".md")

			if !fileExists(srcFile) {
				result.Status = "error"
				result.Detail = fmt.Sprintf("%s %q not found in source", entityType, name)
				results = append(results, result)
				continue
			}

			if fileExists(dstFile) && !force {
				result.Status = "skipped"
				result.Detail = "already exists (use --force to overwrite)"
				results = append(results, result)
				continue
			}

			if dryRun {
				result.Status = "would copy"
				result.Detail = fmt.Sprintf("%s → %s", shortenPath(srcFile), shortenPath(dstFile))
				results = append(results, result)
				continue
			}

			if err := os.MkdirAll(dstDir, 0755); err != nil {
				result.Status = "error"
				result.Detail = err.Error()
				results = append(results, result)
				continue
			}

			if err := copyFile(srcFile, dstFile); err != nil {
				result.Status = "error"
				result.Detail = err.Error()
			} else {
				result.Status = "copied"
				result.Detail = fmt.Sprintf("→ %s", shortenPath(dstFile))
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func copyHooks(fromPath, toPath string, names []string, force, dryRun bool) ([]CopyResult, error) {
	srcClaude, _, err := expandPath(fromPath)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}

	dstClaude, dstIsGlobal, err := expandPath(toPath)
	if err != nil {
		expanded := expandHome(toPath)
		abs, absErr := filepath.Abs(expanded)
		if absErr != nil {
			return nil, fmt.Errorf("destination: %w", err)
		}
		dstClaude = filepath.Join(abs, ".claude")
		dstIsGlobal = false
	}

	// Read source settings
	srcSettings, err := readSettings(srcClaude, fromPath)
	if err != nil {
		return nil, err
	}
	srcHooksRaw, ok := srcSettings["hooks"]
	if !ok {
		return nil, fmt.Errorf("no hooks found in %s", fromPath)
	}
	srcHooksMap, ok := srcHooksRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid hooks format in %s", fromPath)
	}

	// Resolve which events to copy
	isAll := len(names) == 1 && names[0] == "all"
	var eventsToCopy []string
	if isAll {
		for event := range srcHooksMap {
			eventsToCopy = append(eventsToCopy, event)
		}
	} else {
		for _, name := range names {
			if _, exists := srcHooksMap[name]; !exists {
				return nil, fmt.Errorf("hook event %q not found in source", name)
			}
			eventsToCopy = append(eventsToCopy, name)
		}
	}

	// Read or create destination settings
	dstSettingsFile := filepath.Join(dstClaude, settingsFile(dstIsGlobal))
	var dstSettings map[string]any
	if data, err := os.ReadFile(dstSettingsFile); err == nil {
		_ = json.Unmarshal(data, &dstSettings)
	}
	if dstSettings == nil {
		dstSettings = make(map[string]any)
	}

	dstHooksMap, _ := dstSettings["hooks"].(map[string]any)
	if dstHooksMap == nil {
		dstHooksMap = make(map[string]any)
	}

	var results []CopyResult
	for _, event := range eventsToCopy {
		result := CopyResult{
			Name: event,
			From: fromPath,
			To:   toPath,
		}

		if _, exists := dstHooksMap[event]; exists && !force {
			result.Status = "skipped"
			result.Detail = "hook event already exists (use --force to overwrite)"
			results = append(results, result)
			continue
		}

		if dryRun {
			result.Status = "would copy"
			result.Detail = fmt.Sprintf("hook event %q", event)
			results = append(results, result)
			continue
		}

		dstHooksMap[event] = srcHooksMap[event]
		result.Status = "copied"
		result.Detail = fmt.Sprintf("hook event %q → %s", event, shortenPath(dstSettingsFile))
		results = append(results, result)
	}

	if !dryRun {
		dstSettings["hooks"] = dstHooksMap
		if err := os.MkdirAll(dstClaude, 0755); err != nil {
			return results, fmt.Errorf("failed to create destination directory: %w", err)
		}
		data, err := json.MarshalIndent(dstSettings, "", "  ")
		if err != nil {
			return results, fmt.Errorf("failed to marshal settings: %w", err)
		}
		if err := os.WriteFile(dstSettingsFile, data, 0644); err != nil {
			return results, fmt.Errorf("failed to write settings: %w", err)
		}
	}

	return results, nil
}

// resolveNames resolves "all" to actual entity names, or validates specific names.
func resolveNames(srcDir string, entityType string, names []string) ([]string, error) {
	isAll := len(names) == 1 && names[0] == "all"
	if !isAll {
		return names, nil
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read source directory %s: %w", srcDir, err)
	}

	var resolved []string
	for _, e := range entries {
		if entityType == "skills" {
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") && !skipDirs[e.Name()] {
				resolved = append(resolved, e.Name())
			}
		} else {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				resolved = append(resolved, strings.TrimSuffix(e.Name(), ".md"))
			}
		}
	}
	return resolved, nil
}

// copyDir recursively copies a directory tree.
func copyDir(src, dst string) error {
	// Remove destination first if it exists
	if dirExists(dst) {
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("failed to remove existing %s: %w", dst, err)
		}
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		return copyFile(path, target)
	})
}

// copyFile copies a single file, preserving permissions.
func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	abs, _ := filepath.Abs(path)
	if strings.HasPrefix(abs, home) {
		return "~" + abs[len(home):]
	}
	return path
}
