package main

import (
	"testing"

	"github.com/thieung/ccinspect/internal/model"
)

func TestCollectSkills_GlobalAndProject(t *testing.T) {
	inv := &model.Inventory{
		Global: &model.GlobalConfig{
			Skills: []model.Skill{{Name: "g1", Source: "global"}},
		},
		Projects: []*model.ProjectConfig{
			{Skills: []model.Skill{{Name: "p1", Source: "proj"}}},
		},
	}

	all := collectSkills(inv, false)
	if len(all) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(all))
	}

	globalOnly := collectSkills(inv, true)
	if len(globalOnly) != 1 || globalOnly[0].Name != "g1" {
		t.Errorf("global-only = %v", globalOnly)
	}
}

func TestCollectSkills_NilGlobal(t *testing.T) {
	inv := &model.Inventory{
		Projects: []*model.ProjectConfig{
			{Skills: []model.Skill{{Name: "p1"}}},
		},
	}
	all := collectSkills(inv, false)
	if len(all) != 1 {
		t.Errorf("expected 1, got %d", len(all))
	}
}

func TestCollectHooks_GlobalAndProject(t *testing.T) {
	inv := &model.Inventory{
		Global: &model.GlobalConfig{
			Hooks: []model.Hook{{Event: "Pre", Command: "a"}},
		},
		Projects: []*model.ProjectConfig{
			{Hooks: []model.Hook{{Event: "Post", Command: "b"}}},
		},
	}

	all := collectHooks(inv, false)
	if len(all) != 2 {
		t.Fatalf("expected 2, got %d", len(all))
	}

	globalOnly := collectHooks(inv, true)
	if len(globalOnly) != 1 || globalOnly[0].Command != "a" {
		t.Errorf("global-only = %v", globalOnly)
	}
}

func TestCollectHooks_NilGlobal(t *testing.T) {
	inv := &model.Inventory{
		Projects: []*model.ProjectConfig{
			{Hooks: []model.Hook{{Event: "E", Command: "c"}}},
		},
	}
	all := collectHooks(inv, false)
	if len(all) != 1 {
		t.Errorf("expected 1, got %d", len(all))
	}
}

func TestCollectEntities_AllTypes(t *testing.T) {
	inv := &model.Inventory{
		Global: &model.GlobalConfig{
			Agents:   []model.Entity{{Name: "a1", Type: "agent"}},
			Commands: []model.Entity{{Name: "c1", Type: "command"}},
			Rules:    []model.Entity{{Name: "r1", Type: "rule"}},
			Teams:    []model.Entity{{Name: "t1", Type: "team"}},
		},
		Projects: []*model.ProjectConfig{
			{Commands: []model.Entity{{Name: "c2", Type: "command"}}},
		},
	}

	tests := []struct {
		entityType string
		globalOnly bool
		wantCount  int
	}{
		{"agent", false, 1},
		{"agent", true, 1},
		{"command", false, 2},
		{"command", true, 1},
		{"rule", false, 1},
		{"team", false, 1},
	}

	for _, tt := range tests {
		got := collectEntities(inv, tt.entityType, tt.globalOnly)
		if len(got) != tt.wantCount {
			t.Errorf("collectEntities(%q, globalOnly=%v) = %d, want %d", tt.entityType, tt.globalOnly, len(got), tt.wantCount)
		}
	}
}

func TestCollectEntities_NilGlobal(t *testing.T) {
	inv := &model.Inventory{
		Projects: []*model.ProjectConfig{
			{Commands: []model.Entity{{Name: "c1"}}},
		},
	}
	got := collectEntities(inv, "command", false)
	if len(got) != 1 {
		t.Errorf("expected 1, got %d", len(got))
	}
}

func TestCollectMCP(t *testing.T) {
	inv := &model.Inventory{
		Projects: []*model.ProjectConfig{
			{MCPServers: []model.MCPServer{{Name: "s1", Command: "node"}}},
			{MCPServers: []model.MCPServer{{Name: "s2", Command: "python"}}},
		},
	}
	got := collectMCP(inv)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestCollectMCP_Empty(t *testing.T) {
	inv := &model.Inventory{Projects: []*model.ProjectConfig{{}}}
	got := collectMCP(inv)
	if len(got) != 0 {
		t.Errorf("expected 0, got %d", len(got))
	}
}

func TestListCmd_UnknownEntity(t *testing.T) {
	err := listCmd.RunE(listCmd, []string{"unknown-entity"})
	if err == nil {
		t.Error("expected error for unknown entity type")
	}
}
