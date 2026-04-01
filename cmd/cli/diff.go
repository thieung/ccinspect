package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thieunv/ccinspect/internal/config"
	"github.com/thieunv/ccinspect/internal/model"
	"github.com/thieunv/ccinspect/internal/output"
	"github.com/thieunv/ccinspect/internal/parser"
	"github.com/thieunv/ccinspect/internal/scanner"
)

var diffCmd = &cobra.Command{
	Use:   "diff skills <projectA> <projectB>",
	Short: "Compare skills between two projects (use 'global' for ~/.claude)",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		entityType := args[0]
		if entityType != "skills" {
			return fmt.Errorf("only 'skills' diff is supported")
		}

		pathA, pathB := args[1], args[2]
		jsonFlag, _ := cmd.Flags().GetBool("json")

		skillsA := resolveSkills(pathA)
		skillsB := resolveSkills(pathB)

		if jsonFlag {
			fmt.Println(output.RenderJSON(map[string]any{
				"left":  skillsA,
				"right": skillsB,
			}))
		} else {
			fmt.Println(output.RenderDiff(skillsA, skillsB, pathA, pathB))
		}
		return nil
	},
}

func resolveSkills(path string) []model.Skill {
	if path == "global" {
		globalPath, _ := scanner.FindGlobal()
		if globalPath == "" {
			return nil
		}
		return parser.ParseSkills(globalPath, "global")
	}
	cfg := config.Load()
	dirs := scanner.FindClaudeDirs([]string{path}, 1, cfg.ExcludePaths)
	if len(dirs) == 0 {
		return nil
	}
	return parser.ParseSkills(dirs[0], path)
}

func init() {
	diffCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(diffCmd)
}
