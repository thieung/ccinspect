package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// defaultExcludes are directories to skip during scanning.
var defaultExcludes = map[string]bool{
	"node_modules": true,
	".git":         true,
	"vendor":       true,
	".venv":        true,
	"__pycache__":  true,
	"dist":         true,
	"build":        true,
}

// FindGlobal returns the global ~/.claude/ path if it exists.
func FindGlobal() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(home, ".claude")
	if info, err := os.Stat(p); err == nil && info.IsDir() {
		return p, nil
	}
	return "", nil
}

// FindClaudeDirs scans the given paths for directories containing .claude/.
// maxDepth limits how deep to recurse. excludes are directory names to skip.
func FindClaudeDirs(scanPaths []string, maxDepth int, excludes []string) []string {
	excl := make(map[string]bool)
	for k, v := range defaultExcludes {
		excl[k] = v
	}
	for _, e := range excludes {
		excl[e] = true
	}

	var results []string
	seen := make(map[string]bool)

	for _, root := range scanPaths {
		root = expandHome(root)
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // skip inaccessible
			}
			if !d.IsDir() {
				return nil
			}

			// Calculate depth relative to root
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			depth := 0
			if rel != "." {
				depth = strings.Count(rel, string(filepath.Separator)) + 1
			}
			if depth > maxDepth {
				return fs.SkipDir
			}

			name := d.Name()

			// Skip excluded dirs
			if excl[name] && name != filepath.Base(root) {
				return fs.SkipDir
			}

			// Check if this dir contains .claude/
			claudePath := filepath.Join(path, ".claude")
			if info, err := os.Stat(claudePath); err == nil && info.IsDir() {
				abs, _ := filepath.Abs(claudePath)
				if !seen[abs] {
					seen[abs] = true
					results = append(results, abs)
				}
				// Continue recursing — nested projects may also have .claude/
			}

			return nil
		})
	}
	return results
}

func expandHome(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return home
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~"+string(filepath.Separator)) {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
