package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/thieung/ccinspect/internal/copier"
	"github.com/thieung/ccinspect/internal/output"
)

var copyCmd = &cobra.Command{
	Use:   "copy <entity-type> <name> [name2 ...]",
	Short: "Copy skills/agents/hooks/commands between projects or global config",
	Long: `Copy Claude Code entities between projects or between global (~/.claude) and project configs.

Entity types: skills, agents, commands, hooks

Defaults: --from=global, --to=. (current directory)

Sources/destinations:
  "global"         → ~/.claude/ (global config)
  "."              → current directory's .claude/ (project config)
  "/path/to/proj"  → /path/to/proj/.claude/ (project config)

Use "all" as name to copy all entities of that type.

Examples:
  ccinspect copy skills my-skill                                          # global → current dir
  ccinspect copy skills all                                               # all skills from global → current dir
  ccinspect copy agents planner debugger                                  # global agents → current dir
  ccinspect copy skills my-skill --from global --to ~/projects/myapp
  ccinspect copy skills all --from ~/projects/a --to ~/projects/b
  ccinspect copy hooks PreToolUse --from global --to ~/projects/myapp
  ccinspect copy commands my-cmd --from ~/projects/a --to ~/projects/b --force`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		entityType := args[0]
		names := args[1:]
		fromFlag, _ := cmd.Flags().GetString("from")
		toFlag, _ := cmd.Flags().GetString("to")
		forceFlag, _ := cmd.Flags().GetBool("force")
		dryRunFlag, _ := cmd.Flags().GetBool("dry-run")
		jsonFlag, _ := cmd.Flags().GetBool("json")
		listFlag, _ := cmd.Flags().GetBool("list")

		// Validate entity type
		validTypes := map[string]bool{
			"skills": true, "agents": true, "commands": true, "hooks": true,
		}
		if !validTypes[entityType] {
			return fmt.Errorf("unknown entity type: %s (use: skills, agents, commands, hooks)", entityType)
		}

		// List mode: show available entities in the source
		if listFlag {
			if fromFlag == "" {
				return fmt.Errorf("--from is required with --list")
			}
			available, err := copier.ListAvailable(fromFlag, entityType)
			if err != nil {
				return err
			}
			if jsonFlag {
				fmt.Println(output.RenderJSON(available))
			} else {
				titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
				fmt.Println(titleStyle.Render(fmt.Sprintf("Available %s in %s:", entityType, fromFlag)))
				for _, name := range available {
					fmt.Printf("  • %s\n", name)
				}
			}
			return nil
		}

		// Defaults are set via flag defaults: from=global, to=.

		// Perform the copy
		results, err := copier.CopyEntities(fromFlag, toFlag, entityType, names, forceFlag, dryRunFlag)
		if err != nil {
			return err
		}

		if jsonFlag {
			fmt.Println(output.RenderJSON(results))
			return nil
		}

		// Render results
		renderCopyResults(results, dryRunFlag)
		return nil
	},
}

func renderCopyResults(results []copier.CopyResult, dryRun bool) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	skipStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	if dryRun {
		fmt.Println(titleStyle.Render("Dry Run — no files modified"))
	} else {
		fmt.Println(titleStyle.Render("Copy Results"))
	}
	fmt.Println()

	copied, skipped, errors := 0, 0, 0
	for _, r := range results {
		var icon, styledStatus string
		switch r.Status {
		case "copied":
			icon = "✓"
			styledStatus = successStyle.Render(r.Status)
			copied++
		case "would copy":
			icon = "○"
			styledStatus = dimStyle.Render(r.Status)
			copied++
		case "skipped":
			icon = "⊘"
			styledStatus = skipStyle.Render(r.Status)
			skipped++
		default:
			icon = "✗"
			styledStatus = errStyle.Render(r.Status)
			errors++
		}

		detail := ""
		if r.Detail != "" {
			detail = dimStyle.Render(" " + r.Detail)
		}
		fmt.Printf("  %s %s [%s]%s\n", icon, r.Name, styledStatus, detail)
	}

	fmt.Println()

	parts := []string{}
	if copied > 0 {
		label := "copied"
		if dryRun {
			label = "would copy"
		}
		parts = append(parts, successStyle.Render(fmt.Sprintf("%d %s", copied, label)))
	}
	if skipped > 0 {
		parts = append(parts, skipStyle.Render(fmt.Sprintf("%d skipped", skipped)))
	}
	if errors > 0 {
		parts = append(parts, errStyle.Render(fmt.Sprintf("%d errors", errors)))
	}
	fmt.Printf("Summary: %s\n", strings.Join(parts, ", "))
}

func init() {
	copyCmd.Flags().String("from", "global", "Source: 'global' or project path (default: global)")
	copyCmd.Flags().String("to", ".", "Destination: 'global' or project path (default: current directory)")
	copyCmd.Flags().Bool("force", false, "Overwrite existing entities at destination")
	copyCmd.Flags().Bool("dry-run", false, "Show what would be copied without making changes")
	copyCmd.Flags().Bool("json", false, "Output as JSON")
	copyCmd.Flags().Bool("list", false, "List available entities in the source (use with --from)")
	rootCmd.AddCommand(copyCmd)
}
