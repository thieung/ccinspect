package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/thieung/ccinspect/internal/model"
)

// ParseEntities reads .md files from a subdirectory and returns entities.
func ParseEntities(claudePath string, subdir string, entityType string) []model.Entity {
	dir := filepath.Join(claudePath, subdir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var entities []model.Entity
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		entities = append(entities, model.Entity{
			Name: name,
			Path: filepath.Join(dir, e.Name()),
			Type: entityType,
		})
	}
	return entities
}

// ParseTeams reads team directories (each with config.json).
func ParseTeams(claudePath string) []model.Entity {
	dir := filepath.Join(claudePath, "teams")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var teams []model.Entity
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		// Verify config.json exists
		cfg := filepath.Join(dir, e.Name(), "config.json")
		if _, err := os.Stat(cfg); err != nil {
			continue
		}
		teams = append(teams, model.Entity{
			Name: e.Name(),
			Path: filepath.Join(dir, e.Name()),
			Type: "team",
		})
	}
	return teams
}
