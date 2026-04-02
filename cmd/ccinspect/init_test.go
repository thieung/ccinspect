package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCmd_CreatesConfig(t *testing.T) {
	// Save and restore HOME to avoid writing to real home
	origHome := os.Getenv("HOME")
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	cmd := initCmd
	err := cmd.RunE(cmd, nil)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	configFile := filepath.Join(tmp, ".ccinspect", "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("config file was not created")
	}
}
