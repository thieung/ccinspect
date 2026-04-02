package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/thieung/ccinspect/internal/model"
)

// --- helpers ---

func tmpDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// --- ParseSkills ---

func TestParseSkills_WithFrontmatter(t *testing.T) {
	claude := tmpDir(t)
	skillDir := filepath.Join(claude, "skills", "my-skill")
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
description: "A cool skill for testing"
---
# My Skill
Body text here.
`)
	skills := ParseSkills(claude, "global")
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	s := skills[0]
	if s.Name != "my-skill" {
		t.Errorf("name = %q, want my-skill", s.Name)
	}
	if s.Description != "A cool skill for testing" {
		t.Errorf("description = %q", s.Description)
	}
	if s.Source != "global" {
		t.Errorf("source = %q", s.Source)
	}
}

func TestParseSkills_FallbackDescription(t *testing.T) {
	claude := tmpDir(t)
	skillDir := filepath.Join(claude, "skills", "plain")
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `# Title
First content line used as fallback.
`)
	skills := ParseSkills(claude, "proj")
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Description != "First content line used as fallback." {
		t.Errorf("description = %q", skills[0].Description)
	}
}

func TestParseSkills_SkipsDotAndVenv(t *testing.T) {
	claude := tmpDir(t)
	for _, name := range []string{".hidden", ".venv", "__pycache__", "_shared", "legit"} {
		writeFile(t, filepath.Join(claude, "skills", name, "SKILL.md"), "x")
	}
	skills := ParseSkills(claude, "g")
	if len(skills) != 1 || skills[0].Name != "legit" {
		t.Errorf("expected only 'legit', got %v", skills)
	}
}

func TestParseSkills_NoDir(t *testing.T) {
	skills := ParseSkills("/nonexistent", "g")
	if skills != nil {
		t.Errorf("expected nil, got %v", skills)
	}
}

func TestParseSkills_LongDescription(t *testing.T) {
	claude := tmpDir(t)
	long := "description: " + string(make([]byte, 100)) // will be 100 null bytes but let's use real text
	skillDir := filepath.Join(claude, "skills", "long")
	desc := "This is a very long description that exceeds eighty characters and should be truncated by the parser logic"
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), "---\ndescription: "+desc+"\n---\n")
	skills := ParseSkills(claude, "g")
	_ = long
	if len(skills) != 1 {
		t.Fatalf("expected 1, got %d", len(skills))
	}
	if len(skills[0].Description) > 84 { // 80 + "..."
		t.Errorf("description not truncated: len=%d", len(skills[0].Description))
	}
}

func TestParseSkills_WithNameAndPrefix(t *testing.T) {
	claude := tmpDir(t)
	skillDir := filepath.Join(claude, "skills", "claudekit-changelog-sync")
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: ck:changelog-sync
description: "Auto-detect ClaudeKit changelog"
---
# ClaudeKit Changelog Sync
`)
	skills := ParseSkills(claude, "global")
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	s := skills[0]
	if s.Name != "claudekit-changelog-sync" {
		t.Errorf("name = %q, want claudekit-changelog-sync", s.Name)
	}
	if s.DisplayName != "ck:changelog-sync" {
		t.Errorf("display_name = %q, want ck:changelog-sync", s.DisplayName)
	}
	if s.Prefix != "ck" {
		t.Errorf("prefix = %q, want ck", s.Prefix)
	}
	if s.Description != "Auto-detect ClaudeKit changelog" {
		t.Errorf("description = %q", s.Description)
	}
}

func TestParseSkills_NameWithoutPrefix(t *testing.T) {
	claude := tmpDir(t)
	skillDir := filepath.Join(claude, "skills", "my-tool")
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
name: my-tool
description: "A simple tool"
---
`)
	skills := ParseSkills(claude, "global")
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	s := skills[0]
	if s.DisplayName != "my-tool" {
		t.Errorf("display_name = %q, want my-tool", s.DisplayName)
	}
	if s.Prefix != "" {
		t.Errorf("prefix = %q, want empty", s.Prefix)
	}
}

func TestParseSkills_NoNameField(t *testing.T) {
	claude := tmpDir(t)
	skillDir := filepath.Join(claude, "skills", "unnamed")
	writeFile(t, filepath.Join(skillDir, "SKILL.md"), `---
description: "No name provided"
---
`)
	skills := ParseSkills(claude, "global")
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	s := skills[0]
	if s.DisplayName != "" {
		t.Errorf("display_name = %q, want empty", s.DisplayName)
	}
	if s.Prefix != "" {
		t.Errorf("prefix = %q, want empty", s.Prefix)
	}
}

// --- FilterSkillsByPrefix ---

func TestFilterSkillsByPrefix_Matching(t *testing.T) {
	skills := []model.Skill{
		{Name: "claudekit-changelog-sync", DisplayName: "ck:changelog-sync", Prefix: "ck"},
		{Name: "claudekit-docs", DisplayName: "ck:docs", Prefix: "ck"},
		{Name: "skill-debug", DisplayName: "skill:debug", Prefix: "skill"},
		{Name: "plain-skill", Prefix: ""},
	}
	filtered := FilterSkillsByPrefix(skills, "ck")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 skills with prefix 'ck', got %d", len(filtered))
	}
	for _, s := range filtered {
		if s.Prefix != "ck" {
			t.Errorf("unexpected prefix %q in filtered result", s.Prefix)
		}
	}
}

func TestFilterSkillsByPrefix_CaseInsensitive(t *testing.T) {
	skills := []model.Skill{
		{Name: "ck-sync", Prefix: "ck"},
	}
	filtered := FilterSkillsByPrefix(skills, "CK")
	if len(filtered) != 1 {
		t.Errorf("expected case-insensitive match, got %d", len(filtered))
	}
}

func TestFilterSkillsByPrefix_Empty(t *testing.T) {
	skills := []model.Skill{
		{Name: "a", Prefix: "ck"},
		{Name: "b", Prefix: ""},
	}
	filtered := FilterSkillsByPrefix(skills, "")
	if len(filtered) != 2 {
		t.Errorf("empty prefix should return all, got %d", len(filtered))
	}
}

func TestFilterSkillsByPrefix_NoMatch(t *testing.T) {
	skills := []model.Skill{
		{Name: "a", Prefix: "ck"},
	}
	filtered := FilterSkillsByPrefix(skills, "xyz")
	if len(filtered) != 0 {
		t.Errorf("expected 0, got %d", len(filtered))
	}
}

// --- ParseHooks ---

func TestParseHooks_BasicHook(t *testing.T) {
	settings := map[string]any{
		"hooks": map[string]any{
			"PreToolUse": []any{
				map[string]any{
					"matcher": "Bash",
					"hooks": []any{
						map[string]any{
							"command": "echo hello",
							"type":    "command",
						},
					},
				},
			},
		},
	}
	hooks := ParseHooks(settings)
	if len(hooks) != 1 {
		t.Fatalf("expected 1 hook, got %d", len(hooks))
	}
	h := hooks[0]
	if h.Event != "PreToolUse" || h.Matcher != "Bash" || h.Command != "echo hello" || h.Type != "command" {
		t.Errorf("hook = %+v", h)
	}
}

func TestParseHooks_DefaultType(t *testing.T) {
	settings := map[string]any{
		"hooks": map[string]any{
			"PostToolUse": []any{
				map[string]any{
					"matcher": "",
					"hooks": []any{
						map[string]any{"command": "cmd"},
					},
				},
			},
		},
	}
	hooks := ParseHooks(settings)
	if len(hooks) != 1 || hooks[0].Type != "command" {
		t.Errorf("expected default type 'command', got %+v", hooks)
	}
}

func TestParseHooks_NilSettings(t *testing.T) {
	if hooks := ParseHooks(nil); hooks != nil {
		t.Errorf("expected nil")
	}
}

func TestParseHooks_NoHooksKey(t *testing.T) {
	if hooks := ParseHooks(map[string]any{"foo": "bar"}); hooks != nil {
		t.Errorf("expected nil")
	}
}

func TestParseHooks_TruncatesLongCommand(t *testing.T) {
	long := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 70 chars
	settings := map[string]any{
		"hooks": map[string]any{
			"E": []any{
				map[string]any{
					"matcher": "M",
					"hooks":   []any{map[string]any{"command": long}},
				},
			},
		},
	}
	hooks := ParseHooks(settings)
	if len(hooks) != 1 {
		t.Fatalf("expected 1 hook")
	}
	if len(hooks[0].Command) > 60 {
		t.Errorf("command not truncated: len=%d", len(hooks[0].Command))
	}
}

// --- ParseEntities ---

func TestParseEntities_MdFiles(t *testing.T) {
	claude := tmpDir(t)
	writeFile(t, filepath.Join(claude, "rules", "alpha.md"), "content")
	writeFile(t, filepath.Join(claude, "rules", "beta.md"), "content")
	writeFile(t, filepath.Join(claude, "rules", "skip.txt"), "not md")
	os.MkdirAll(filepath.Join(claude, "rules", "subdir"), 0o755)

	entities := ParseEntities(claude, "rules", "rule")
	if len(entities) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(entities))
	}
	names := []string{entities[0].Name, entities[1].Name}
	sort.Strings(names)
	if names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("names = %v", names)
	}
	if entities[0].Type != "rule" {
		t.Errorf("type = %q", entities[0].Type)
	}
}

func TestParseEntities_EmptyDir(t *testing.T) {
	claude := tmpDir(t)
	os.MkdirAll(filepath.Join(claude, "agents"), 0o755)
	entities := ParseEntities(claude, "agents", "agent")
	if len(entities) != 0 {
		t.Errorf("expected 0, got %d", len(entities))
	}
}

func TestParseEntities_MissingDir(t *testing.T) {
	entities := ParseEntities("/nonexistent", "x", "y")
	if entities != nil {
		t.Errorf("expected nil")
	}
}

// --- ParseTeams ---

func TestParseTeams_WithConfig(t *testing.T) {
	claude := tmpDir(t)
	writeFile(t, filepath.Join(claude, "teams", "alpha", "config.json"), "{}")
	// team without config.json should be skipped
	os.MkdirAll(filepath.Join(claude, "teams", "no-config"), 0o755)

	teams := ParseTeams(claude)
	if len(teams) != 1 || teams[0].Name != "alpha" {
		t.Errorf("teams = %v", teams)
	}
}

// --- ParseSettings ---

func TestParseSettings_Valid(t *testing.T) {
	claude := tmpDir(t)
	writeFile(t, filepath.Join(claude, "settings.json"), `{"key":"val"}`)
	m, err := ParseSettings(claude, "settings.json")
	if err != nil {
		t.Fatal(err)
	}
	if m["key"] != "val" {
		t.Errorf("got %v", m)
	}
}

func TestParseSettings_Missing(t *testing.T) {
	_, err := ParseSettings("/nonexistent", "settings.json")
	if err == nil {
		t.Error("expected error")
	}
}

func TestParseSettings_InvalidJSON(t *testing.T) {
	claude := tmpDir(t)
	writeFile(t, filepath.Join(claude, "s.json"), `{bad`)
	_, err := ParseSettings(claude, "s.json")
	if err == nil {
		t.Error("expected error")
	}
}

// --- ParseMCP ---

func TestParseMCP_Valid(t *testing.T) {
	// ParseMCP expects .mcp.json at parent of claudePath
	projectRoot := tmpDir(t)
	claudePath := filepath.Join(projectRoot, ".claude")
	os.MkdirAll(claudePath, 0o755)

	mcpData := map[string]any{
		"mcpServers": map[string]any{
			"my-server": map[string]any{
				"command": "node",
				"args":    []string{"server.js", "--port=3000"},
			},
		},
	}
	data, _ := json.Marshal(mcpData)
	writeFile(t, filepath.Join(projectRoot, ".mcp.json"), string(data))

	servers := ParseMCP(claudePath)
	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}
	if servers[0].Name != "my-server" || servers[0].Command != "node" {
		t.Errorf("server = %+v", servers[0])
	}
	if len(servers[0].Args) != 2 {
		t.Errorf("args = %v", servers[0].Args)
	}
}

func TestParseMCP_Missing(t *testing.T) {
	servers := ParseMCP("/nonexistent/.claude")
	if servers != nil {
		t.Errorf("expected nil")
	}
}

// --- BuildInventory ---

func TestBuildInventory_Integration(t *testing.T) {
	// Set up a fake global claude path
	globalRoot := tmpDir(t)
	globalClaude := filepath.Join(globalRoot, ".claude")

	// Global skills
	writeFile(t, filepath.Join(globalClaude, "skills", "sk1", "SKILL.md"), "---\ndescription: skill one\n---\n")
	// Global rules
	writeFile(t, filepath.Join(globalClaude, "rules", "my-rule.md"), "rule content")
	// Global settings with hooks
	settings := map[string]any{
		"hooks": map[string]any{
			"PreToolUse": []any{
				map[string]any{
					"matcher": "Bash",
					"hooks":   []any{map[string]any{"command": "check"}},
				},
			},
		},
	}
	settingsData, _ := json.Marshal(settings)
	writeFile(t, filepath.Join(globalClaude, "settings.json"), string(settingsData))
	// Global CLAUDE.md
	writeFile(t, filepath.Join(globalClaude, "CLAUDE.md"), "# Global")

	// Set up a fake project
	projectRoot := tmpDir(t)
	projectClaude := filepath.Join(projectRoot, ".claude")
	writeFile(t, filepath.Join(projectClaude, "settings.local.json"), `{}`)
	writeFile(t, filepath.Join(projectRoot, "CLAUDE.md"), "# Project")
	// Project MCP
	mcpData := map[string]any{
		"mcpServers": map[string]any{
			"srv": map[string]any{"command": "cmd", "args": []string{}},
		},
	}
	mcpBytes, _ := json.Marshal(mcpData)
	writeFile(t, filepath.Join(projectRoot, ".mcp.json"), string(mcpBytes))

	inv := BuildInventory(globalClaude, []string{globalClaude, projectClaude})

	// Verify global
	if inv.Global == nil {
		t.Fatal("global is nil")
	}
	if len(inv.Global.Skills) != 1 {
		t.Errorf("global skills = %d", len(inv.Global.Skills))
	}
	if len(inv.Global.Rules) != 1 {
		t.Errorf("global rules = %d", len(inv.Global.Rules))
	}
	if len(inv.Global.Hooks) != 1 {
		t.Errorf("global hooks = %d", len(inv.Global.Hooks))
	}
	if !inv.Global.HasClaudeMD {
		t.Error("expected HasClaudeMD true")
	}

	// Verify project (global path should be skipped)
	if len(inv.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(inv.Projects))
	}
	proj := inv.Projects[0]
	if !proj.HasClaudeMD {
		t.Error("project HasClaudeMD should be true")
	}
	if len(proj.MCPServers) != 1 {
		t.Errorf("project mcp servers = %d", len(proj.MCPServers))
	}
}

func TestBuildInventory_EmptyGlobal(t *testing.T) {
	inv := BuildInventory("", nil)
	if inv.Global != nil {
		t.Error("expected nil global")
	}
	if len(inv.Projects) != 0 {
		t.Error("expected no projects")
	}
}
