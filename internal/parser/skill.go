package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/thieung/ccinspect/internal/model"
)

// skipDirs are skill subdirectories to ignore.
var skipDirs = map[string]bool{
	"_shared": true, ".venv": true, "__pycache__": true,
}

// ParseSkills scans a skills/ directory and returns skill metadata.
func ParseSkills(claudePath string, source string) []model.Skill {
	skillsDir := filepath.Join(claudePath, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil
	}

	var skills []model.Skill
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || skipDirs[e.Name()] {
			continue
		}
		skill := model.Skill{
			Name:   e.Name(),
			Path:   filepath.Join(skillsDir, e.Name()),
			Source: source,
		}
		// Extract name + description from SKILL.md frontmatter
		meta := extractSkillMeta(filepath.Join(skill.Path, "SKILL.md"))
		skill.Description = meta.description
		if meta.name != "" {
			skill.DisplayName = meta.name
			if idx := strings.Index(meta.name, ":"); idx > 0 {
				skill.Prefix = meta.name[:idx]
			}
		}
		skills = append(skills, skill)
	}
	return skills
}

// FilterSkillsByPrefix filters skills to only those matching the given prefix.
func FilterSkillsByPrefix(skills []model.Skill, prefix string) []model.Skill {
	if prefix == "" {
		return skills
	}
	var filtered []model.Skill
	for _, s := range skills {
		if strings.EqualFold(s.Prefix, prefix) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

type skillMeta struct {
	name        string
	description string
}

func extractSkillMeta(path string) skillMeta {
	f, err := os.Open(path)
	if err != nil {
		return skillMeta{}
	}
	defer f.Close()

	var meta skillMeta
	sc := bufio.NewScanner(f)
	inFrontmatter := false
	lines := 0
	for sc.Scan() && lines < 20 {
		line := strings.TrimSpace(sc.Text())
		lines++
		if line == "---" {
			inFrontmatter = !inFrontmatter
			continue
		}
		if inFrontmatter {
			// Extract name from frontmatter
			if strings.HasPrefix(line, "name:") {
				name := strings.TrimPrefix(line, "name:")
				name = strings.TrimSpace(name)
				name = strings.Trim(name, `"'`)
				meta.name = name
			}
			// Extract description from frontmatter
			if strings.HasPrefix(line, "description:") {
				desc := strings.TrimPrefix(line, "description:")
				desc = strings.TrimSpace(desc)
				desc = strings.Trim(desc, `"'`)
				if len(desc) > 80 {
					meta.description = desc[:80] + "..."
				} else {
					meta.description = desc
				}
			}
			continue
		}
		if meta.description != "" {
			continue
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Fallback: first content line after frontmatter
		if len(line) > 80 {
			meta.description = line[:80] + "..."
		} else {
			meta.description = line
		}
	}
	return meta
}
