package model

// Skill represents a Claude Code skill (folder with SKILL.md).
type Skill struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"` // from SKILL.md frontmatter name field (e.g. "ck:changelog-sync")
	Prefix      string `json:"prefix,omitempty"`       // prefix extracted from display name (e.g. "ck")
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source"` // "global" or project path
}

// Hook represents a hook entry from settings.json.
type Hook struct {
	Event   string `json:"event"`
	Matcher string `json:"matcher,omitempty"`
	Command string `json:"command"`
	Type    string `json:"type"`
}

// MCPServer represents an MCP server configuration.
type MCPServer struct {
	Name    string   `json:"name"`
	Type    string   `json:"type,omitempty"` // e.g. "http"
	URL     string   `json:"url,omitempty"`  // For HTTP servers
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// Entity represents a generic named entity (agent, command, rule, team).
type Entity struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Type   string `json:"type"` // "agent", "command", "rule", "team"
	Prefix string `json:"prefix,omitempty"` // prefix extracted from name (e.g., "ck" from "ck-planner")
}
