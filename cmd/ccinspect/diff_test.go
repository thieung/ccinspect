package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiffCmd_OnlySkillsSupported(t *testing.T) {
	err := diffCmd.RunE(diffCmd, []string{"hooks", "/tmp", "/tmp"})
	if err == nil {
		t.Error("expected error for non-skills entity type")
	}
}

func TestResolveSkills_NonexistentPath(t *testing.T) {
	skills := resolveSkills("/nonexistent/path")
	if len(skills) != 0 {
		t.Errorf("expected 0 skills for nonexistent path, got %d", len(skills))
	}
}

func TestResolveSkills_WithSkillDir(t *testing.T) {
	tmp := t.TempDir()
	claudeDir := filepath.Join(tmp, ".claude")
	skillDir := filepath.Join(claudeDir, "skills", "test-skill")
	os.MkdirAll(skillDir, 0o755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("---\ndescription: test\n---\n"), 0o644)

	skills := resolveSkills(tmp)
	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}
}
