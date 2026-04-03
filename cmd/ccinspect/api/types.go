package api

// ScanResponse is the response for GET /api/scan.
type ScanResponse struct {
	Projects []ProjectSummary `json:"projects"`
	Global  *GlobalSummary   `json:"global,omitempty"`
}

// ProjectSummary is a short description of a scanned project.
type ProjectSummary struct {
	Path       string `json:"path"`
	Skills     int    `json:"skills"`
	Hooks      int    `json:"hooks"`
	Agents     int    `json:"agents"`
	Commands   int    `json:"commands"`
	MCPServers int    `json:"mcp_servers"`
	Teams      int    `json:"teams"`
	HasClaudeMD bool  `json:"has_claude_md"`
}

// GlobalSummary is a short description of the global config.
type GlobalSummary struct {
	Path        string `json:"path"`
	Skills      int    `json:"skills"`
	Hooks       int    `json:"hooks"`
	Agents      int    `json:"agents"`
	Commands    int    `json:"commands"`
	Rules       int    `json:"rules"`
	Teams       int    `json:"teams"`
	HasClaudeMD bool   `json:"has_claude_md"`
}

// ListResponse is the response for GET /api/list.
type ListResponse struct {
	Entities []EntityItem `json:"entities"`
}

// EntityItem is a single entity returned by the list endpoint.
type EntityItem struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source"`
	Type        string `json:"type,omitempty"`
}

// DiffResponse is the response for GET /api/diff.
type DiffResponse struct {
	LeftName  string      `json:"left_name"`
	RightName string      `json:"right_name"`
	EntityType string     `json:"entity_type"`
	Added     []DiffEntry `json:"added"`
	Removed   []DiffEntry `json:"removed"`
	Changed   []DiffEntry `json:"changed"`
}

// DiffEntry represents a single diff line.
type DiffEntry struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Side        string `json:"side,omitempty"` // "left", "right", "both"
}

// CopyRequest is the body for POST /api/copy.
type CopyRequest struct {
	Type  string   `json:"type"`
	Names []string `json:"names"`
	From  string   `json:"from"`
	To    string   `json:"to"`
	Force bool     `json:"force"`
	DryRun bool    `json:"dry_run"`
}

// CopyResponse is the response for POST /api/copy.
type CopyResponse struct {
	Status  string       `json:"status"`
	Results []CopyResult `json:"results,omitempty"`
	Message string       `json:"message,omitempty"`
}

// CopyResult mirrors copier.CopyResult.
type CopyResult struct {
	Name   string `json:"name"`
	From   string `json:"from"`
	To     string `json:"to"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// CleanRequest is the body for POST /api/clean.
type CleanRequest struct {
	Path   string `json:"path"`
	DryRun bool   `json:"dry_run"`
}

// CleanResponse is the response for POST /api/clean.
type CleanResponse struct {
	Status      string   `json:"status"`
	Message     string   `json:"message,omitempty"`
	FilesCount  int      `json:"files_count"`
	Files       []string `json:"files,omitempty"`
	DryRun      bool     `json:"dry_run"`
}

// CleanTeamsRequest is the body for POST /api/clean/teams.
type CleanTeamsRequest struct {
	TeamNames []string `json:"team_names,omitempty"` // empty = all
	All       bool     `json:"all"`
	DryRun    bool     `json:"dry_run"`
}

// CleanTeamsResponse is the response for POST /api/clean/teams.
type CleanTeamsResponse struct {
	Status       string   `json:"status"`
	Message      string   `json:"message,omitempty"`
	TeamsCount   int      `json:"teams_count"`
	TeamNames    []string `json:"team_names,omitempty"`
	DryRun       bool     `json:"dry_run"`
}

// ConfigResponse is the response for GET /api/config.
type ConfigResponse struct {
	Config ConfigData `json:"config"`
}

// ConfigData mirrors config.Config.
type ConfigData struct {
	ScanPaths     []string `json:"scan_paths"`
	ExcludePaths  []string `json:"exclude_paths"`
	MaxDepth      int      `json:"max_depth"`
	DefaultOutput string   `json:"default_output"`
}

// SaveConfigRequest is the body for POST /api/config.
type SaveConfigRequest struct {
	ScanPaths     []string `json:"scan_paths"`
	ExcludePaths  []string `json:"exclude_paths"`
	MaxDepth      int      `json:"max_depth"`
	DefaultOutput string   `json:"default_output"`
}

// ErrorResponse is a generic API error.
type ErrorResponse struct {
	Error string `json:"error"`
}
