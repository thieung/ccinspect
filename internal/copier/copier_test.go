package copier

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// setupTestDirs creates a source and destination .claude directory structure for testing.
func setupTestDirs(t *testing.T) (srcRoot, dstRoot string) {
	t.Helper()
	srcRoot = t.TempDir()
	dstRoot = t.TempDir()

	srcClaude := filepath.Join(srcRoot, ".claude")

	// Create skills
	for _, skill := range []string{"my-skill", "another-skill"} {
		skillDir := filepath.Join(srcClaude, "skills", skill)
		os.MkdirAll(skillDir, 0755)
		os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("---\ndescription: test skill\n---\nTest content"), 0644)
	}
	// Create _shared (should be skipped)
	os.MkdirAll(filepath.Join(srcClaude, "skills", "_shared"), 0755)

	// Create agents
	agentsDir := filepath.Join(srcClaude, "agents")
	os.MkdirAll(agentsDir, 0755)
	os.WriteFile(filepath.Join(agentsDir, "test-agent.md"), []byte("# Test Agent"), 0644)

	// Create commands
	cmdsDir := filepath.Join(srcClaude, "commands")
	os.MkdirAll(cmdsDir, 0755)
	os.WriteFile(filepath.Join(cmdsDir, "test-cmd.md"), []byte("# Test Command"), 0644)

	// Create settings.json with hooks
	settings := map[string]any{
		"hooks": map[string]any{
			"PreToolUse": []any{
				map[string]any{
					"matcher": "",
					"hooks": []any{
						map[string]any{
							"type":    "command",
							"command": "echo pre-tool",
						},
					},
				},
			},
			"PostToolUse": []any{
				map[string]any{
					"matcher": "",
					"hooks": []any{
						map[string]any{
							"type":    "command",
							"command": "echo post-tool",
						},
					},
				},
			},
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	os.WriteFile(filepath.Join(srcClaude, "settings.json"), data, 0644)

	// Create destination .claude (empty)
	os.MkdirAll(filepath.Join(dstRoot, ".claude"), 0755)

	return srcRoot, dstRoot
}

func TestListAvailable_Skills(t *testing.T) {
	srcRoot, _ := setupTestDirs(t)

	names, err := ListAvailable(srcRoot, "skills")
	if err != nil {
		t.Fatalf("ListAvailable failed: %v", err)
	}

	// Should contain my-skill and another-skill but NOT _shared
	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	if !nameSet["my-skill"] {
		t.Error("expected my-skill in list")
	}
	if !nameSet["another-skill"] {
		t.Error("expected another-skill in list")
	}
	if nameSet["_shared"] {
		t.Error("_shared should be filtered out")
	}
}

func TestListAvailable_Agents(t *testing.T) {
	srcRoot, _ := setupTestDirs(t)

	names, err := ListAvailable(srcRoot, "agents")
	if err != nil {
		t.Fatalf("ListAvailable failed: %v", err)
	}

	if len(names) != 1 || names[0] != "test-agent" {
		t.Errorf("expected [test-agent], got %v", names)
	}
}

func TestListAvailable_Commands(t *testing.T) {
	srcRoot, _ := setupTestDirs(t)

	names, err := ListAvailable(srcRoot, "commands")
	if err != nil {
		t.Fatalf("ListAvailable failed: %v", err)
	}

	if len(names) != 1 || names[0] != "test-cmd" {
		t.Errorf("expected [test-cmd], got %v", names)
	}
}

func TestListAvailable_Hooks(t *testing.T) {
	srcRoot, _ := setupTestDirs(t)

	events, err := ListAvailable(srcRoot, "hooks")
	if err != nil {
		t.Fatalf("ListAvailable failed: %v", err)
	}

	eventSet := make(map[string]bool)
	for _, e := range events {
		eventSet[e] = true
	}

	if !eventSet["PreToolUse"] || !eventSet["PostToolUse"] {
		t.Errorf("expected PreToolUse and PostToolUse, got %v", events)
	}
}

func TestCopySkills_Single(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	results, err := CopyEntities(srcRoot, dstRoot, "skills", []string{"my-skill"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "copied" {
		t.Errorf("expected status 'copied', got %q", results[0].Status)
	}

	// Verify the skill was actually copied
	copiedSkillMD := filepath.Join(dstRoot, ".claude", "skills", "my-skill", "SKILL.md")
	if _, err := os.Stat(copiedSkillMD); err != nil {
		t.Errorf("expected copied SKILL.md at %s", copiedSkillMD)
	}
}

func TestCopySkills_All(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	results, err := CopyEntities(srcRoot, dstRoot, "skills", []string{"all"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results (excluding _shared), got %d", len(results))
	}
	for _, r := range results {
		if r.Status != "copied" {
			t.Errorf("expected status 'copied' for %s, got %q", r.Name, r.Status)
		}
		if r.Name == "_shared" {
			t.Error("_shared should not be copied")
		}
	}
}

func TestCopySkills_SkipExisting(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	// Copy once
	CopyEntities(srcRoot, dstRoot, "skills", []string{"my-skill"}, false, false)

	// Copy again without force — should be skipped
	results, err := CopyEntities(srcRoot, dstRoot, "skills", []string{"my-skill"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if results[0].Status != "skipped" {
		t.Errorf("expected status 'skipped', got %q", results[0].Status)
	}
}

func TestCopySkills_ForceOverwrite(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	// Copy once
	CopyEntities(srcRoot, dstRoot, "skills", []string{"my-skill"}, false, false)

	// Copy again with force — should succeed
	results, err := CopyEntities(srcRoot, dstRoot, "skills", []string{"my-skill"}, true, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if results[0].Status != "copied" {
		t.Errorf("expected status 'copied' with --force, got %q", results[0].Status)
	}
}

func TestCopySkills_DryRun(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	results, err := CopyEntities(srcRoot, dstRoot, "skills", []string{"my-skill"}, false, true)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if results[0].Status != "would copy" {
		t.Errorf("expected status 'would copy', got %q", results[0].Status)
	}

	// Verify nothing was actually copied
	copiedSkillMD := filepath.Join(dstRoot, ".claude", "skills", "my-skill", "SKILL.md")
	if _, err := os.Stat(copiedSkillMD); err == nil {
		t.Error("dry run should not create files")
	}
}

func TestCopySkills_NotFound(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	results, err := CopyEntities(srcRoot, dstRoot, "skills", []string{"nonexistent-skill"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if results[0].Status != "error" {
		t.Errorf("expected status 'error' for missing skill, got %q", results[0].Status)
	}
}

func TestCopyAgents_Single(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	results, err := CopyEntities(srcRoot, dstRoot, "agents", []string{"test-agent"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if results[0].Status != "copied" {
		t.Errorf("expected status 'copied', got %q", results[0].Status)
	}

	copiedAgent := filepath.Join(dstRoot, ".claude", "agents", "test-agent.md")
	data, err := os.ReadFile(copiedAgent)
	if err != nil {
		t.Fatalf("expected copied agent at %s", copiedAgent)
	}
	if string(data) != "# Test Agent" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestCopyCommands_Single(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	results, err := CopyEntities(srcRoot, dstRoot, "commands", []string{"test-cmd"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if results[0].Status != "copied" {
		t.Errorf("expected status 'copied', got %q", results[0].Status)
	}

	copiedCmd := filepath.Join(dstRoot, ".claude", "commands", "test-cmd.md")
	if _, err := os.Stat(copiedCmd); err != nil {
		t.Errorf("expected copied command at %s", copiedCmd)
	}
}

func TestCopyHooks_Single(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	results, err := CopyEntities(srcRoot, dstRoot, "hooks", []string{"PreToolUse"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if len(results) != 1 || results[0].Status != "copied" {
		t.Errorf("expected 1 copied result, got %v", results)
	}

	// Verify settings file was created with hook
	settingsFile := filepath.Join(dstRoot, ".claude", "settings.local.json")
	data, err := os.ReadFile(settingsFile)
	if err != nil {
		t.Fatalf("expected settings file at %s", settingsFile)
	}

	var settings map[string]any
	json.Unmarshal(data, &settings)
	hooks, ok := settings["hooks"].(map[string]any)
	if !ok {
		t.Fatal("expected hooks in settings")
	}
	if _, exists := hooks["PreToolUse"]; !exists {
		t.Error("expected PreToolUse hook in destination settings")
	}
}

func TestCopyHooks_All(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	results, err := CopyEntities(srcRoot, dstRoot, "hooks", []string{"all"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 hook events copied, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != "copied" {
			t.Errorf("expected status 'copied' for %s, got %q", r.Name, r.Status)
		}
	}
}

func TestCopyHooks_SkipExisting(t *testing.T) {
	srcRoot, dstRoot := setupTestDirs(t)

	// Copy once
	CopyEntities(srcRoot, dstRoot, "hooks", []string{"PreToolUse"}, false, false)

	// Copy again without force — should be skipped
	results, err := CopyEntities(srcRoot, dstRoot, "hooks", []string{"PreToolUse"}, false, false)
	if err != nil {
		t.Fatalf("CopyEntities failed: %v", err)
	}

	if results[0].Status != "skipped" {
		t.Errorf("expected status 'skipped', got %q", results[0].Status)
	}
}

func TestCopyHooks_NotFound(t *testing.T) {
	srcRoot, _ := setupTestDirs(t)
	dstRoot := t.TempDir()
	os.MkdirAll(filepath.Join(dstRoot, ".claude"), 0755)

	_, err := CopyEntities(srcRoot, dstRoot, "hooks", []string{"NonExistentEvent"}, false, false)
	if err == nil {
		t.Error("expected error for nonexistent hook event")
	}
}

func TestExpandPath_Global(t *testing.T) {
	path, isGlobal, err := expandPath("global")
	if err != nil {
		t.Fatalf("expandPath failed: %v", err)
	}
	if !isGlobal {
		t.Error("expected isGlobal=true for 'global'")
	}
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".claude")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestExpandPath_ProjectPath(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, ".claude"), 0755)

	path, isGlobal, err := expandPath(tmpDir)
	if err != nil {
		t.Fatalf("expandPath failed: %v", err)
	}
	if isGlobal {
		t.Error("expected isGlobal=false for project path")
	}
	expected := filepath.Join(tmpDir, ".claude")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestExpandPath_NoClaudeDir(t *testing.T) {
	tmpDir := t.TempDir()
	// No .claude directory

	_, _, err := expandPath(tmpDir)
	if err == nil {
		t.Error("expected error when no .claude/ exists")
	}
}
