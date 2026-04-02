package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandCleanPath_Dot(t *testing.T) {
	got := expandCleanPath(".")
	wd, _ := os.Getwd()
	if got != wd {
		t.Errorf("expandCleanPath(\".\") = %q, want %q", got, wd)
	}
}

func TestExpandCleanPath_Tilde(t *testing.T) {
	got := expandCleanPath("~/projects")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, "projects")
	if got != want {
		t.Errorf("expandCleanPath(\"~/projects\") = %q, want %q", got, want)
	}
}

func TestExpandCleanPath_Absolute(t *testing.T) {
	got := expandCleanPath("/tmp/test")
	if got != "/tmp/test" {
		t.Errorf("expandCleanPath(\"/tmp/test\") = %q", got)
	}
}

func TestExpandCleanPath_Relative(t *testing.T) {
	got := expandCleanPath("some/path")
	if !filepath.IsAbs(got) {
		t.Errorf("expected absolute path, got %q", got)
	}
}

func TestCleanCmd_NoClaude(t *testing.T) {
	tmp := t.TempDir()
	err := cleanCmd.RunE(cleanCmd, []string{tmp})
	if err == nil {
		t.Error("expected error when no .claude/ exists")
	}
}

func TestCleanCmd_DryRun(t *testing.T) {
	tmp := t.TempDir()
	claudeDir := filepath.Join(tmp, ".claude")
	os.MkdirAll(claudeDir, 0o755)
	os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte("{}"), 0o644)

	// Set dry-run flag
	cleanCmd.ResetFlags()
	cleanCmd.Flags().Bool("dry-run", false, "Show what would be removed without deleting")
	cleanCmd.Flags().Set("dry-run", "true")

	err := cleanCmd.RunE(cleanCmd, []string{tmp})
	if err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}

	// .claude/ should still exist after dry-run
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		t.Error(".claude/ was deleted during dry-run")
	}
}

func TestCleanCmd_ActualClean(t *testing.T) {
	tmp := t.TempDir()
	claudeDir := filepath.Join(tmp, ".claude")
	os.MkdirAll(claudeDir, 0o755)
	os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte("{}"), 0o644)
	os.WriteFile(filepath.Join(tmp, ".mcp.json"), []byte("{}"), 0o644)

	// Reset flags to ensure dry-run is false
	cleanCmd.ResetFlags()
	cleanCmd.Flags().Bool("dry-run", false, "Show what would be removed without deleting")

	err := cleanCmd.RunE(cleanCmd, []string{tmp})
	if err != nil {
		t.Fatalf("clean failed: %v", err)
	}

	if _, err := os.Stat(claudeDir); !os.IsNotExist(err) {
		t.Error(".claude/ should be removed")
	}
	if _, err := os.Stat(filepath.Join(tmp, ".mcp.json")); !os.IsNotExist(err) {
		t.Error(".mcp.json should be removed")
	}
}

func setupTeamsDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	teamsDir := filepath.Join(tmp, ".claude", "teams")

	// Stale team (no config.json)
	os.MkdirAll(filepath.Join(teamsDir, "stale-uuid-team"), 0o755)

	// Configured team (has config.json)
	configuredDir := filepath.Join(teamsDir, "my-configured-team")
	os.MkdirAll(configuredDir, 0o755)
	os.WriteFile(filepath.Join(configuredDir, "config.json"), []byte("{}"), 0o644)

	return tmp
}

func TestCleanTeamsCmd_StaleOnly_DryRun(t *testing.T) {
	tmp := setupTeamsDir(t)
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	cleanTeamsCmd.ResetFlags()
	cleanTeamsCmd.Flags().Bool("dry-run", false, "")
	cleanTeamsCmd.Flags().Bool("all", false, "")
	cleanTeamsCmd.Flags().Set("dry-run", "true")

	err := cleanTeamsCmd.RunE(cleanTeamsCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both dirs should still exist (dry-run)
	teamsDir := filepath.Join(tmp, ".claude", "teams")
	if _, err := os.Stat(filepath.Join(teamsDir, "stale-uuid-team")); os.IsNotExist(err) {
		t.Error("stale team should not be deleted in dry-run")
	}
}

func TestCleanTeamsCmd_StaleOnly(t *testing.T) {
	tmp := setupTeamsDir(t)
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	cleanTeamsCmd.ResetFlags()
	cleanTeamsCmd.Flags().Bool("dry-run", false, "")
	cleanTeamsCmd.Flags().Bool("all", false, "")

	err := cleanTeamsCmd.RunE(cleanTeamsCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	teamsDir := filepath.Join(tmp, ".claude", "teams")
	if _, err := os.Stat(filepath.Join(teamsDir, "stale-uuid-team")); !os.IsNotExist(err) {
		t.Error("stale team should be removed")
	}
	if _, err := os.Stat(filepath.Join(teamsDir, "my-configured-team")); os.IsNotExist(err) {
		t.Error("configured team should be kept")
	}
}

func TestCleanTeamsCmd_All(t *testing.T) {
	tmp := setupTeamsDir(t)
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	cleanTeamsCmd.ResetFlags()
	cleanTeamsCmd.Flags().Bool("dry-run", false, "")
	cleanTeamsCmd.Flags().Bool("all", false, "")
	cleanTeamsCmd.Flags().Set("all", "true")

	err := cleanTeamsCmd.RunE(cleanTeamsCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	teamsDir := filepath.Join(tmp, ".claude", "teams")
	if _, err := os.Stat(filepath.Join(teamsDir, "stale-uuid-team")); !os.IsNotExist(err) {
		t.Error("stale team should be removed")
	}
	if _, err := os.Stat(filepath.Join(teamsDir, "my-configured-team")); !os.IsNotExist(err) {
		t.Error("configured team should also be removed with --all")
	}
}
