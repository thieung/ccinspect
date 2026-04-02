package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thieung/ccinspect/internal/output"
)

var cleanCmd = &cobra.Command{
	Use:   "clean <project-path>",
	Short: "Remove .claude/ directory from a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		target := expandCleanPath(args[0])
		claudeDir := filepath.Join(target, ".claude")

		info, err := os.Stat(claudeDir)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("no .claude/ directory found at %s", target)
		}

		// Show what will be removed
		spin := output.NewSpinner("Analyzing .claude/ contents...")
		spin.Start()

		var files []string
		filepath.Walk(claudeDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		spin.Stop()

		fmt.Printf("Found %d files in %s\n", len(files), claudeDir)

		if dryRun {
			for _, f := range files {
				rel, _ := filepath.Rel(target, f)
				fmt.Printf("  %s\n", rel)
			}
			fmt.Println("\nDry run — no files removed. Run without --dry-run to delete.")
			return nil
		}

		if err := os.RemoveAll(claudeDir); err != nil {
			return fmt.Errorf("failed to remove %s: %w", claudeDir, err)
		}

		// Also remove .mcp.json if present
		mcpFile := filepath.Join(target, ".mcp.json")
		if _, err := os.Stat(mcpFile); err == nil {
			os.Remove(mcpFile)
			fmt.Println("Removed .mcp.json")
		}

		fmt.Printf("Cleaned .claude/ from %s (%d files removed)\n", target, len(files))
		return nil
	},
}

var cleanTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Remove stale or all teams from ~/.claude/teams/",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		all, _ := cmd.Flags().GetBool("all")

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory: %w", err)
		}
		teamsDir := filepath.Join(home, ".claude", "teams")

		entries, err := os.ReadDir(teamsDir)
		if err != nil {
			return fmt.Errorf("no teams directory found at %s", teamsDir)
		}

		var toRemove []string
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			teamPath := filepath.Join(teamsDir, e.Name())
			hasConfig := false
			if _, err := os.Stat(filepath.Join(teamPath, "config.json")); err == nil {
				hasConfig = true
			}

			if all || !hasConfig {
				label := "stale"
				if hasConfig {
					label = "configured"
				}
				toRemove = append(toRemove, teamPath)
				if dryRun {
					fmt.Printf("  [%s] %s\n", label, e.Name())
				}
			}
		}

		if len(toRemove) == 0 {
			fmt.Println("No teams to clean.")
			return nil
		}

		if dryRun {
			fmt.Printf("\nDry run — %d teams would be removed. Run without --dry-run to delete.\n", len(toRemove))
			return nil
		}

		for _, p := range toRemove {
			if err := os.RemoveAll(p); err != nil {
				fmt.Fprintf(os.Stderr, "failed to remove %s: %v\n", p, err)
			}
		}
		fmt.Printf("Removed %d teams from %s\n", len(toRemove), teamsDir)
		return nil
	},
}

func init() {
	cleanCmd.Flags().Bool("dry-run", false, "Show what would be removed without deleting")
	cleanTeamsCmd.Flags().Bool("dry-run", false, "Show what would be removed without deleting")
	cleanTeamsCmd.Flags().Bool("all", false, "Remove all teams (not just stale ones)")
	cleanCmd.AddCommand(cleanTeamsCmd)
	rootCmd.AddCommand(cleanCmd)
}

func expandCleanPath(path string) string {
	if path == "." {
		abs, _ := os.Getwd()
		return abs
	}
	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
	if len(path) > 1 && (path[:2] == "~/" || path[:2] == "~"+string(filepath.Separator)) {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
