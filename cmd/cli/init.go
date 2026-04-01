package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thieunv/ccinspect/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default config at ~/.ccinspect/config.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.DefaultConfig()
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("Config created at ~/.ccinspect/config.json")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
