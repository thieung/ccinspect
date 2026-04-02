package output

import (
	"fmt"
	"strings"

	"github.com/thieung/ccinspect/internal/model"
)

// RenderInventoryMarkdown prints the main scan summary as a markdown table.
func RenderInventoryMarkdown(inv *model.Inventory, titleText string, hideEmpty bool) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s\n\n", titleText))
	sb.WriteString("| Project | Skills | Hooks | Agents | MCP |\n")
	sb.WriteString("|---------|--------|-------|--------|-----|\n")

	totalSkills, totalHooks, totalAgents, totalMCP := 0, 0, 0, 0

	// Global row
	if inv.Global != nil {
		g := inv.Global
		s, h, a := len(g.Skills), len(g.Hooks), len(g.Agents)
		if !hideEmpty || s+h+a > 0 {
			totalSkills += s
			totalHooks += h
			totalAgents += a
			sb.WriteString(fmt.Sprintf("| `~/.claude` (GLOBAL) | %d | %d | %d | - |\n", s, h, a))
		}
	}

	displayedProjects := 0
	for _, p := range inv.Projects {
		s, h, a, m := len(p.Skills), len(p.Hooks), len(p.Agents), len(p.MCPServers)
		if hideEmpty && s+h+a+m == 0 {
			continue
		}
		displayedProjects++
		totalSkills += s
		totalHooks += h
		totalAgents += a
		totalMCP += m
		sb.WriteString(fmt.Sprintf("| `%s` | %d | %d | %d | %d |\n", shortenPath(p.Path), s, h, a, m))
	}

	sb.WriteString(fmt.Sprintf("| **TOTAL: %d projects** | **%d** | **%d** | **%d** | **%d** |\n", displayedProjects, totalSkills, totalHooks, totalAgents, totalMCP))

	return sb.String()
}

// RenderEntityListMarkdown prints entities as a markdown list or table.
func RenderEntityListMarkdown(entities []model.Entity, header string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## %s (%d)\n\n", header, len(entities)))
	if len(entities) == 0 {
		sb.WriteString("*(none)*\n")
		return sb.String()
	}

	sb.WriteString("| Name | Path |\n")
	sb.WriteString("|------|------|\n")
	for _, e := range entities {
		sb.WriteString(fmt.Sprintf("| **%s** | `%s` |\n", e.Name, shortenPath(e.Path)))
	}

	return sb.String()
}

// RenderSkillListMarkdown prints skills as a markdown table.
func RenderSkillListMarkdown(skills []model.Skill, header string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## %s (%d)\n\n", header, len(skills)))
	if len(skills) == 0 {
		sb.WriteString("*(none)*\n")
		return sb.String()
	}

	hasDisplayName := false
	for _, s := range skills {
		if s.DisplayName != "" {
			hasDisplayName = true
			break
		}
	}

	if hasDisplayName {
		sb.WriteString("| Display Name | Dir Name | Description | Source |\n")
		sb.WriteString("|--------------|----------|-------------|--------|\n")
		for _, s := range skills {
			display := s.DisplayName
			if display == "" {
				display = s.Name
			}
			desc := strings.ReplaceAll(s.Description, "\n", " ")
			sb.WriteString(fmt.Sprintf("| **%s** | `%s` | %s | %s |\n", display, s.Name, desc, s.Source))
		}
	} else {
		sb.WriteString("| Name | Description | Source |\n")
		sb.WriteString("|------|-------------|--------|\n")
		for _, s := range skills {
			desc := strings.ReplaceAll(s.Description, "\n", " ")
			sb.WriteString(fmt.Sprintf("| **%s** | %s | %s |\n", s.Name, desc, s.Source))
		}
	}

	return sb.String()
}

// RenderHookListMarkdown prints hooks as a markdown table.
func RenderHookListMarkdown(hooks []model.Hook, header string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## %s (%d)\n\n", header, len(hooks)))
	if len(hooks) == 0 {
		sb.WriteString("*(none)*\n")
		return sb.String()
	}

	sb.WriteString("| Event | Matcher | Command |\n")
	sb.WriteString("|-------|---------|---------|\n")
	for _, h := range hooks {
		sb.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` |\n", h.Event, h.Matcher, h.Command))
	}

	return sb.String()
}

// RenderMCPListMarkdown prints MCP servers as a markdown table.
func RenderMCPListMarkdown(servers []model.MCPServer, header string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## %s (%d)\n\n", header, len(servers)))
	if len(servers) == 0 {
		sb.WriteString("*(none)*\n")
		return sb.String()
	}

	sb.WriteString("| Name | Command / Type | Args / URL |\n")
	sb.WriteString("|------|----------------|------------|\n")
	for _, s := range servers {
		if s.Type == "http" || s.URL != "" {
			sb.WriteString(fmt.Sprintf("| **%s** | `http` | %s |\n", s.Name, s.URL))
		} else {
			sb.WriteString(fmt.Sprintf("| **%s** | `%s` | `%s` |\n", s.Name, s.Command, strings.Join(s.Args, " ")))
		}
	}

	return sb.String()
}

// RenderDiffMarkdown prints a skill diff as a markdown table.
func RenderDiffMarkdown(left, right []model.Skill, leftName, rightName string) string {
	var sb strings.Builder
	sb.WriteString("## Skill Diff\n\n")

	leftMap := make(map[string]bool)
	rightMap := make(map[string]bool)
	for _, s := range left {
		leftMap[s.Name] = true
	}
	for _, s := range right {
		rightMap[s.Name] = true
	}

	sb.WriteString(fmt.Sprintf("| Skill | %s | %s |\n", shortenPath(leftName), shortenPath(rightName)))
	sb.WriteString(fmt.Sprintf("|-------|%s|%s|\n", strings.Repeat("-", len(shortenPath(leftName))), strings.Repeat("-", len(shortenPath(rightName)))))

	// Only in left
	for _, s := range left {
		if !rightMap[s.Name] {
			sb.WriteString(fmt.Sprintf("| **%s** | ✓ | - |\n", s.Name))
		}
	}
	// Only in right
	for _, s := range right {
		if !leftMap[s.Name] {
			sb.WriteString(fmt.Sprintf("| **%s** | - | ✓ |\n", s.Name))
		}
	}
	// In both
	for _, s := range left {
		if rightMap[s.Name] {
			sb.WriteString(fmt.Sprintf("| **%s** | ✓ | ✓ |\n", s.Name))
		}
	}

	return sb.String()
}
