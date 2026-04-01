package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/thieunv/ccinspect/internal/model"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).MarginBottom(1)
)

// RenderInventoryTable prints the main scan summary table.
func RenderInventoryTable(inv *model.Inventory) string {
	rows := [][]string{}
	totalSkills, totalHooks, totalAgents, totalMCP := 0, 0, 0, 0

	// Global row
	if inv.Global != nil {
		g := inv.Global
		s, h, a := len(g.Skills), len(g.Hooks), len(g.Agents)
		totalSkills += s
		totalHooks += h
		totalAgents += a
		rows = append(rows, []string{
			"~/.claude (GLOBAL)", itoa(s), itoa(h), itoa(a), "-",
		})
	}

	// Project rows
	for _, p := range inv.Projects {
		s, h, m := len(p.Skills), len(p.Hooks), len(p.MCPServers)
		totalSkills += s
		totalHooks += h
		totalMCP += m
		rows = append(rows, []string{
			shortenPath(p.Path), itoa(s), itoa(h), "-", itoa(m),
		})
	}

	// Total row
	rows = append(rows, []string{
		fmt.Sprintf("TOTAL: %d projects", len(inv.Projects)),
		itoa(totalSkills), itoa(totalHooks), itoa(totalAgents), itoa(totalMCP),
	})

	t := table.New().
		Headers("Project", "Skills", "Hooks", "Agents", "MCP").
		Rows(rows...).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			// Last row (total) is bold
			if row == len(rows)-1 {
				return lipgloss.NewStyle().Bold(true)
			}
			return lipgloss.NewStyle()
		})

	title := titleStyle.Render("Claude Code Installations")
	return fmt.Sprintf("%s\n%s", title, t.String())
}

// RenderEntityList prints a list of entities with name and optional info.
func RenderEntityList(entities []model.Entity, header string) string {
	if len(entities) == 0 {
		return fmt.Sprintf("%s: (none)\n", header)
	}

	rows := [][]string{}
	for _, e := range entities {
		rows = append(rows, []string{e.Name, e.Path})
	}

	t := table.New().
		Headers("Name", "Path").
		Rows(rows...).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return lipgloss.NewStyle()
		})

	title := titleStyle.Render(fmt.Sprintf("%s (%d)", header, len(entities)))
	return fmt.Sprintf("%s\n%s", title, t.String())
}

// RenderSkillList prints skills with description.
func RenderSkillList(skills []model.Skill, header string) string {
	if len(skills) == 0 {
		return fmt.Sprintf("%s: (none)\n", header)
	}

	rows := [][]string{}
	for _, s := range skills {
		rows = append(rows, []string{s.Name, s.Description, s.Source})
	}

	t := table.New().
		Headers("Name", "Description", "Source").
		Rows(rows...).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return lipgloss.NewStyle()
		})

	title := titleStyle.Render(fmt.Sprintf("%s (%d)", header, len(skills)))
	return fmt.Sprintf("%s\n%s", title, t.String())
}

// RenderHookList prints hooks with event/matcher/command.
func RenderHookList(hooks []model.Hook, header string) string {
	if len(hooks) == 0 {
		return fmt.Sprintf("%s: (none)\n", header)
	}

	rows := [][]string{}
	for _, h := range hooks {
		rows = append(rows, []string{h.Event, h.Matcher, h.Command})
	}

	t := table.New().
		Headers("Event", "Matcher", "Command").
		Rows(rows...).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return lipgloss.NewStyle()
		})

	title := titleStyle.Render(fmt.Sprintf("%s (%d)", header, len(hooks)))
	return fmt.Sprintf("%s\n%s", title, t.String())
}

// RenderMCPList prints MCP servers.
func RenderMCPList(servers []model.MCPServer, header string) string {
	if len(servers) == 0 {
		return fmt.Sprintf("%s: (none)\n", header)
	}

	rows := [][]string{}
	for _, s := range servers {
		rows = append(rows, []string{s.Name, s.Command, strings.Join(s.Args, " ")})
	}

	t := table.New().
		Headers("Name", "Command", "Args").
		Rows(rows...).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return lipgloss.NewStyle()
		})

	title := titleStyle.Render(fmt.Sprintf("%s (%d)", header, len(servers)))
	return fmt.Sprintf("%s\n%s", title, t.String())
}

// RenderDiff prints side-by-side skill comparison.
func RenderDiff(left, right []model.Skill, leftName, rightName string) string {
	leftMap := make(map[string]bool)
	rightMap := make(map[string]bool)
	for _, s := range left {
		leftMap[s.Name] = true
	}
	for _, s := range right {
		rightMap[s.Name] = true
	}

	rows := [][]string{}
	// Only in left
	for _, s := range left {
		if !rightMap[s.Name] {
			rows = append(rows, []string{s.Name, "✓", "-"})
		}
	}
	// Only in right
	for _, s := range right {
		if !leftMap[s.Name] {
			rows = append(rows, []string{s.Name, "-", "✓"})
		}
	}
	// In both
	for _, s := range left {
		if rightMap[s.Name] {
			rows = append(rows, []string{s.Name, "✓", "✓"})
		}
	}

	t := table.New().
		Headers("Skill", shortenPath(leftName), shortenPath(rightName)).
		Rows(rows...).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return lipgloss.NewStyle()
		})

	title := titleStyle.Render("Skill Diff")
	return fmt.Sprintf("%s\n%s", title, t.String())
}

func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err == nil {
		abs, _ := filepath.Abs(path)
		if strings.HasPrefix(abs, home) {
			return "~" + abs[len(home):]
		}
	}
	parts := strings.Split(path, string(filepath.Separator))
	if len(parts) > 3 {
		return "~/" + strings.Join(parts[len(parts)-2:], "/")
	}
	return path
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
