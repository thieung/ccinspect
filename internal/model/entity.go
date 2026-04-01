package model

// Skill represents a Claude Code skill (folder with SKILL.md).
type Skill struct {
	Name        string `json:"name"`
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

// MCPServer represents an MCP server from .mcp.json.
type MCPServer struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
}

// Entity represents a generic named entity (agent, command, rule, team).
type Entity struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"` // "agent", "command", "rule", "team"
}
