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
		mdFlag, _ := cmd.Flags().GetBool("md")
		globalOnly, _ := cmd.Flags().GetBool("global-only")
		projectFlag, _ := cmd.Flags().GetString("project")
		prefixFlag, _ := cmd.Flags().GetString("prefix")

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
			if prefixFlag != "" {
				skills = parser.FilterSkillsByPrefix(skills, prefixFlag)
			}
			header := "Skills"
			if prefixFlag != "" {
				header = fmt.Sprintf("Skills (prefix: %s)", prefixFlag)
			}
			if jsonFlag {
				fmt.Println(output.RenderJSON(skills))
			} else if mdFlag {
				fmt.Print(output.RenderSkillListMarkdown(skills, header))
			} else {
				fmt.Println(output.RenderSkillList(skills, header))
			}
		case "hooks":
			hooks := collectHooks(inv, globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(hooks))
			} else if mdFlag {
				fmt.Print(output.RenderHookListMarkdown(hooks, "Hooks"))
			} else {
				fmt.Println(output.RenderHookList(hooks, "Hooks"))
			}
		case "agents":
			entities := collectEntities(inv, "agent", globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(entities))
			} else if mdFlag {
				fmt.Print(output.RenderEntityListMarkdown(entities, "Agents"))
			} else {
				fmt.Println(output.RenderEntityList(entities, "Agents"))
			}
		case "commands":
			entities := collectEntities(inv, "command", globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(entities))
			} else if mdFlag {
				fmt.Print(output.RenderEntityListMarkdown(entities, "Commands"))
			} else {
				fmt.Println(output.RenderEntityList(entities, "Commands"))
			}
		case "rules":
			entities := collectEntities(inv, "rule", globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(entities))
			} else if mdFlag {
				fmt.Print(output.RenderEntityListMarkdown(entities, "Rules"))
			} else {
				fmt.Println(output.RenderEntityList(entities, "Rules"))
			}
		case "teams":
			entities := collectEntities(inv, "team", globalOnly)
			if jsonFlag {
				fmt.Println(output.RenderJSON(entities))
			} else if mdFlag {
				fmt.Print(output.RenderEntityListMarkdown(entities, "Teams"))
			} else {
				fmt.Println(output.RenderEntityList(entities, "Teams"))
			}
		case "mcp":
			servers := collectMCP(inv)
			if jsonFlag {
				fmt.Println(output.RenderJSON(servers))
			} else if mdFlag {
				fmt.Print(output.RenderMCPListMarkdown(servers, "MCP Servers"))
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
			} else if entityType == "agent" {
				all = append(all, p.Agents...)
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
	listCmd.Flags().Bool("md", false, "Output as Markdown")
	listCmd.Flags().Bool("global-only", false, "Show only global entities")
	listCmd.Flags().String("project", "", "Show only entities from a specific project")
	listCmd.Flags().String("prefix", "", "Filter skills by prefix (e.g. ck, skill)")
	rootCmd.AddCommand(listCmd)
}
