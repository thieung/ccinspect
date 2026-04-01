package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thieunv/ccinspect/internal/output"
	"github.com/thieunv/ccinspect/internal/parser"
	"github.com/thieunv/ccinspect/internal/scanner"
)

var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "Show global Claude Code configuration summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonFlag, _ := cmd.Flags().GetBool("json")

		globalPath, _ := scanner.FindGlobal()
		if globalPath == "" {
			return fmt.Errorf("no global ~/.claude/ directory found")
		}

		inv := parser.BuildInventory(globalPath, nil)

		if jsonFlag {
			fmt.Println(output.RenderJSON(inv.Global))
			return nil
		}

		g := inv.Global
		fmt.Printf("Global Config: %s\n\n", g.Path)
		fmt.Printf("  Skills:   %d\n", len(g.Skills))
		fmt.Printf("  Hooks:    %d\n", len(g.Hooks))
		fmt.Printf("  Agents:   %d\n", len(g.Agents))
		fmt.Printf("  Commands: %d\n", len(g.Commands))
		fmt.Printf("  Rules:    %d\n", len(g.Rules))
		fmt.Printf("  Teams:    %d\n", len(g.Teams))
		fmt.Printf("  CLAUDE.md: %v\n", g.HasClaudeMD)
		return nil
	},
}

func init() {
	globalCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(globalCmd)
}
