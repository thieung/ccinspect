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
		// Extract description from SKILL.md first non-empty, non-frontmatter line
		skill.Description = extractSkillDescription(filepath.Join(skill.Path, "SKILL.md"))
		skills = append(skills, skill)
	}
	return skills
}

func extractSkillDescription(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

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
			// Extract description from frontmatter
			if strings.HasPrefix(line, "description:") {
				desc := strings.TrimPrefix(line, "description:")
				desc = strings.TrimSpace(desc)
				desc = strings.Trim(desc, `"'`)
				if len(desc) > 80 {
					return desc[:80] + "..."
				}
				return desc
			}
			continue
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Fallback: first content line after frontmatter
		if len(line) > 80 {
			return line[:80] + "..."
		}
		return line
	}
	return ""
}
