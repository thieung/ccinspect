package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thieung/ccinspect/internal/config"
	"github.com/thieung/ccinspect/internal/output"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		fmt.Println(output.RenderJSON(cfg))
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
