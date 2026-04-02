package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thieung/ccinspect/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default config at ~/.ccinspect/config.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.DefaultConfig()
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		home, _ := os.UserHomeDir()
		fmt.Printf("Config created at %s\n", filepath.Join(home, ".ccinspect", "config.json"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
