package model

import "time"

// Inventory is the top-level result of scanning all Claude configs.
type Inventory struct {
	Global    *GlobalConfig    `json:"global"`
	Projects  []*ProjectConfig `json:"projects"`
	ScannedAt time.Time        `json:"scanned_at"`
}

// GlobalConfig represents ~/.claude/ configuration.
type GlobalConfig struct {
	Path        string         `json:"path"`
	Skills      []Skill        `json:"skills"`
	Hooks       []Hook         `json:"hooks"`
	Agents      []Entity       `json:"agents"`
	Commands    []Entity       `json:"commands"`
	Rules       []Entity       `json:"rules"`
	Teams       []Entity       `json:"teams"`
	Settings    map[string]any `json:"settings,omitempty"`
	HasClaudeMD bool           `json:"has_claude_md"`
}

// ProjectConfig represents a project-level .claude/ configuration.
type ProjectConfig struct {
	Path        string         `json:"path"`
	ClaudePath  string         `json:"claude_path"`
	Skills      []Skill        `json:"skills"`
	Hooks       []Hook         `json:"hooks"`
	Agents      []Entity       `json:"agents"` // Custom agents at project level
	Commands    []Entity       `json:"commands"`
	MCPServers  []MCPServer    `json:"mcp_servers"`
	HasClaudeMD bool           `json:"has_claude_md"`
	Settings    map[string]any `json:"settings,omitempty"`
}
