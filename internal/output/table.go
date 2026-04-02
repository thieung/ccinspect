package output

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/thieung/ccinspect/internal/model"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00E5FF")) // Bright Cyan
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E84A6E")).MarginBottom(1) // Pink/Magenta
	borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Grey border
)

// RenderInventoryTable prints the main scan summary table.
func RenderInventoryTable(inv *model.Inventory, titleText string, hideEmpty bool) string {
	rows := [][]string{}
	totalSkills, totalHooks, totalAgents, totalMCP := 0, 0, 0, 0

	// Global row
	if inv.Global != nil {
		g := inv.Global
		s, h, a := len(g.Skills), len(g.Hooks), len(g.Agents)
		if !hideEmpty || s+h+a > 0 {
			totalSkills += s
			totalHooks += h
			totalAgents += a
			rows = append(rows, []string{
				"~/.claude (GLOBAL)", itoa(s), itoa(h), itoa(a), "-",
			})
		}
	}

	// Sort projects by total entities (skills+hooks+agents+mcp) descending
	sorted := make([]*model.ProjectConfig, len(inv.Projects))
	copy(sorted, inv.Projects)
	sort.Slice(sorted, func(i, j int) bool {
		ti := len(sorted[i].Skills) + len(sorted[i].Hooks) + len(sorted[i].Agents) + len(sorted[i].MCPServers)
		tj := len(sorted[j].Skills) + len(sorted[j].Hooks) + len(sorted[j].Agents) + len(sorted[j].MCPServers)
		return ti > tj
	})

	// Project rows
	displayedProjects := 0
	for _, p := range sorted {
		s, h, a, m := len(p.Skills), len(p.Hooks), len(p.Agents), len(p.MCPServers)
		if hideEmpty && s+h+a+m == 0 {
			continue // Skip empty projects
		}
		displayedProjects++
		totalSkills += s
		totalHooks += h
		totalAgents += a
		totalMCP += m
		rows = append(rows, []string{
			shortenPath(p.Path), itoa(s), itoa(h), itoa(a), itoa(m),
		})
	}

	// Total row
	rows = append(rows, []string{
		fmt.Sprintf("TOTAL: %d projects", displayedProjects),
		itoa(totalSkills), itoa(totalHooks), itoa(totalAgents), itoa(totalMCP),
	})

	t := table.New().
		Headers("Project", "Skills", "Hooks", "Agents", "MCP").
		Rows(rows...).
		Border(lipgloss.RoundedBorder()).BorderStyle(borderStyle).
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

	title := titleStyle.Render(titleText)
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
		Border(lipgloss.RoundedBorder()).BorderStyle(borderStyle).
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

	// Check if any skill has a display name (prefixed skills)
	hasDisplayName := false
	for _, s := range skills {
		if s.DisplayName != "" {
			hasDisplayName = true
			break
		}
	}

	rows := [][]string{}
	var headers []string
	if hasDisplayName {
		headers = []string{"Display Name", "Dir Name", "Description", "Source"}
		for _, s := range skills {
			display := s.DisplayName
			if display == "" {
				display = s.Name
			}
			rows = append(rows, []string{display, s.Name, s.Description, s.Source})
		}
	} else {
		headers = []string{"Name", "Description", "Source"}
		for _, s := range skills {
			rows = append(rows, []string{s.Name, s.Description, s.Source})
		}
	}

	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.RoundedBorder()).BorderStyle(borderStyle).
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
		Border(lipgloss.RoundedBorder()).BorderStyle(borderStyle).
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
		if s.Type == "http" || s.URL != "" {
			rows = append(rows, []string{s.Name, "http", s.URL})
		} else {
			rows = append(rows, []string{s.Name, s.Command, strings.Join(s.Args, " ")})
		}
	}

	t := table.New().
		Headers("Name", "Command / Type", "Args / URL").
		Rows(rows...).
		Border(lipgloss.RoundedBorder()).BorderStyle(borderStyle).
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
		Border(lipgloss.RoundedBorder()).BorderStyle(borderStyle).
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
			rel := abs[len(home):]
			return "~" + filepath.ToSlash(rel)
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
