package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thieung/ccinspect/internal/config"
	"github.com/thieung/ccinspect/internal/model"
	"github.com/thieung/ccinspect/internal/output"
	"github.com/thieung/ccinspect/internal/parser"
	"github.com/thieung/ccinspect/internal/scanner"
)

var listCmd = &cobra.Command{
	Use:   "list <entity>",
	Short: "List all entities of a type (skills, hooks, agents, commands, rules, mcp, teams)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entityType := args[0]
		jsonFlag, _ := cmd.Flags().GetBool("json")
		globalOnly, _ := cmd.Flags().GetBool("global-only")
		projectFlag, _ := cmd.Flags().GetString("project")

		cfg := config.Load()
		globalPath, _ := scanner.FindGlobal()

		var claudePaths []string
		if !globalOnly {
			if projectFlag != "" {
				claudePaths = scanner.FindClaudeDirs([]string{projectFlag}, 1, cfg.ExcludePaths)
			} else {
				claudePaths = scanner.FindClaudeDirs(cfg.ScanPaths, cfg.MaxDepth, cfg.ExcludePaths)
			}
		}

		inv := parser.BuildInventory(globalPath, claudePaths)

		switch entityType {
		case "skills":
			skills := collectSkills(inv, globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(skills))
			} else {
				fmt.Println(output.RenderSkillList(skills, "Skills"))
			}
		case "hooks":
			hooks := collectHooks(inv, globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(hooks))
			} else {
				fmt.Println(output.RenderHookList(hooks, "Hooks"))
			}
		case "agents":
			entities := collectEntities(inv, "agent", globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(entities))
			} else {
				fmt.Println(output.RenderEntityList(entities, "Agents"))
			}
		case "commands":
			entities := collectEntities(inv, "command", globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(entities))
			} else {
				fmt.Println(output.RenderEntityList(entities, "Commands"))
			}
		case "rules":
			entities := collectEntities(inv, "rule", globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(entities))
			} else {
				fmt.Println(output.RenderEntityList(entities, "Rules"))
			}
		case "teams":
			entities := collectEntities(inv, "team", globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(entities))
			} else {
				fmt.Println(output.RenderEntityList(entities, "Teams"))
			}
		case "mcp":
			servers := collectMCP(inv)
			if jsonFlag {
				fmt.Println(output.RenderJSON(servers))
			} else {
				fmt.Println(output.RenderMCPList(servers, "MCP Servers"))
			}
		default:
			return fmt.Errorf("unknown entity type: %s (use: skills, hooks, agents, commands, rules, mcp, teams)", entityType)
		}
		return nil
	},
}

func collectSkills(inv *model.Inventory, globalOnly bool) []model.Skill {
	var all []model.Skill
	if inv.Global != nil {
		all = append(all, inv.Global.Skills...)
	}
	if !globalOnly {
		for _, p := range inv.Projects {
			all = append(all, p.Skills...)
		}
	}
	return all
}

func collectHooks(inv *model.Inventory, globalOnly bool) []model.Hook {
	var all []model.Hook
	if inv.Global != nil {
		all = append(all, inv.Global.Hooks...)
	}
	if !globalOnly {
		for _, p := range inv.Projects {
			all = append(all, p.Hooks...)
		}
	}
	return all
}

func collectEntities(inv *model.Inventory, entityType string, globalOnly bool) []model.Entity {
	var all []model.Entity
	if inv.Global != nil {
		switch entityType {
		case "agent":
			all = append(all, inv.Global.Agents...)
		case "command":
			all = append(all, inv.Global.Commands...)
		case "rule":
			all = append(all, inv.Global.Rules...)
		case "team":
			all = append(all, inv.Global.Teams...)
		}
	}
	if !globalOnly {
		for _, p := range inv.Projects {
			if entityType == "command" {
				all = append(all, p.Commands...)
			}
		}
	}
	return all
}

func collectMCP(inv *model.Inventory) []model.MCPServer {
	var all []model.MCPServer
	for _, p := range inv.Projects {
		all = append(all, p.MCPServers...)
	}
	return all
}

func init() {
	listCmd.Flags().Bool("json", false, "Output as JSON")
	listCmd.Flags().Bool("global-only", false, "Show only global entities")
	listCmd.Flags().String("project", "", "Show only entities from a specific project")
	rootCmd.AddCommand(listCmd)
}
