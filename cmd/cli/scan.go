package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thieung/ccinspect/internal/config"
	"github.com/thieung/ccinspect/internal/output"
	"github.com/thieung/ccinspect/internal/parser"
	"github.com/thieung/ccinspect/internal/scanner"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan all Claude Code installations and show summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		projectFlag, _ := cmd.Flags().GetString("project")

		cfg := config.Load()

		spin := output.NewSpinner("Scanning projects...")
		if !jsonFlag {
			spin.Start()
		}

		// Find global
		globalPath, _ := scanner.FindGlobal()

		var claudePaths []string
		if projectFlag != "" {
			// Single project mode
			claudePaths = scanner.FindClaudeDirs([]string{projectFlag}, 1, cfg.ExcludePaths)
		} else {
			claudePaths = scanner.FindClaudeDirs(cfg.ScanPaths, cfg.MaxDepth, cfg.ExcludePaths)
		}

		inv := parser.BuildInventory(globalPath, claudePaths)

		if !jsonFlag {
			spin.Stop()
		}

		if jsonFlag {
			fmt.Println(output.RenderJSON(inv))
		} else {
			fmt.Println(output.RenderInventoryTable(inv))
		}
		return nil
	},
}

func init() {
	scanCmd.Flags().Bool("json", false, "Output as JSON")
	scanCmd.Flags().String("project", "", "Scan a single project path")
	rootCmd.AddCommand(scanCmd)
}
