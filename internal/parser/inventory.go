package parser

import (
	"os"
	"path/filepath"
	"time"

	"github.com/thieung/ccinspect/internal/model"
)

// BuildInventory scans global and project paths, returning a full inventory.
func BuildInventory(globalPath string, claudePaths []string) *model.Inventory {
	inv := &model.Inventory{
		ScannedAt: time.Now(),
	}

	// Parse global config
	if globalPath != "" {
		inv.Global = parseGlobal(globalPath)
	}

	// Parse each project
	for _, cp := range claudePaths {
		// Skip global path if it appears in project list
		if cp == globalPath {
			continue
		}
		proj := parseProject(cp)
		if proj != nil {
			inv.Projects = append(inv.Projects, proj)
		}
	}

	return inv
}

func parseGlobal(claudePath string) *model.GlobalConfig {
	gc := &model.GlobalConfig{
		Path: claudePath,
	}

	// Settings
	settings, _ := ParseSettings(claudePath, "settings.json")
	gc.Settings = settings

	// Entities
	gc.Skills = ParseSkills(claudePath, "global")
	gc.Hooks = ParseHooks(settings)
	gc.Agents = ParseEntities(claudePath, "agents", "agent")
	gc.Commands = ParseEntities(claudePath, "commands", "command")
	gc.Rules = ParseEntities(claudePath, "rules", "rule")
	gc.Teams = ParseTeams(claudePath)

	// CLAUDE.md check
	gc.HasClaudeMD = fileExists(filepath.Join(claudePath, "CLAUDE.md"))

	return gc
}

func parseProject(claudePath string) *model.ProjectConfig {
	projectRoot := filepath.Dir(claudePath)
	pc := &model.ProjectConfig{
		Path:       projectRoot,
		ClaudePath: claudePath,
	}

	// Settings (project uses settings.local.json)
	settings, _ := ParseSettings(claudePath, "settings.local.json")
	pc.Settings = settings

	// Entities
	pc.Skills = ParseSkills(claudePath, projectRoot)
	pc.Hooks = ParseHooks(settings)
	pc.Commands = ParseEntities(claudePath, "commands", "command")
	pc.MCPServers = ParseMCP(claudePath)

	// CLAUDE.md at project root
	pc.HasClaudeMD = fileExists(filepath.Join(projectRoot, "CLAUDE.md"))

	return pc
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
