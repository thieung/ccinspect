package main

import (
	"testing"
)

func TestConfigShowCmd_Executes(t *testing.T) {
	// config show just prints default config JSON — should not error
	cmd := configShowCmd
	err := cmd.RunE(cmd, nil)
	if err != nil {
		t.Fatalf("config show failed: %v", err)
	}
}
