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
		prefixFlag, _ := cmd.Flags().GetString("prefix")

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

		// Apply prefix filter
		if prefixFlag != "" {
			if inv.Global != nil {
				inv.Global.Skills = parser.FilterSkillsByPrefix(inv.Global.Skills, prefixFlag)
			}
			for _, p := range inv.Projects {
				p.Skills = parser.FilterSkillsByPrefix(p.Skills, prefixFlag)
			}
		}

		if !jsonFlag {
			spin.Stop()
		}

		if jsonFlag {
			fmt.Println(output.RenderJSON(inv))
		} else {
			if prefixFlag != "" {
				fmt.Println(output.RenderInventoryTable(inv, fmt.Sprintf("Claude Code Installations (prefix: %s)", prefixFlag), true))
			} else {
				fmt.Println(output.RenderInventoryTable(inv, "Claude Code Installations", false))
			}
		}
		return nil
	},
}

func init() {
	scanCmd.Flags().Bool("json", false, "Output as JSON")
	scanCmd.Flags().String("project", "", "Scan a single project path")
	scanCmd.Flags().String("prefix", "", "Filter skills by prefix (e.g. ck, skill)")
	rootCmd.AddCommand(scanCmd)
}
