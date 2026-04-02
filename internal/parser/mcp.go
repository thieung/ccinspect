package parser

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/thieung/ccinspect/internal/model"
)

// mcpFileFormat represents a file containing mcpServers at the top level (e.g., .mcp.json).
type mcpFileFormat struct {
	MCPServers map[string]struct {
		Type    string   `json:"type"`
		URL     string   `json:"url"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
	} `json:"mcpServers"`
}

// claudeJSONFormat represents the format of the global ~/.claude.json file.
type claudeJSONFormat struct {
	Projects map[string]struct {
		MCPServers map[string]struct {
			Type    string   `json:"type"`
			URL     string   `json:"url"`
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"mcpServers"`
	} `json:"projects"`
}

// ParseMCP reads MCP server configs from all known locations:
// 1. .mcp.json at project root
// 2. .claude/settings.json (mcpServers key)
// 3. .claude/settings.local.json (mcpServers key)
func ParseMCP(claudePath string) []model.MCPServer {
	projectRoot := filepath.Dir(claudePath)
	seen := make(map[string]bool)
	var servers []model.MCPServer

	// Sources to check, in order
	sources := []string{
		filepath.Join(projectRoot, ".mcp.json"),
		filepath.Join(projectRoot, ".claude.json"), // Support local .claude.json scope
		filepath.Join(claudePath, "settings.json"),
		filepath.Join(claudePath, "settings.local.json"),
	}

	for _, path := range sources {
		for _, srv := range parseMCPFile(path) {
			if !seen[srv.Name] {
				seen[srv.Name] = true
				servers = append(servers, srv)
			}
		}
	}

	// 4. Global ~/.claude.json project registry
	for _, srv := range parseClaudeJSON(projectRoot) {
		if !seen[srv.Name] {
			seen[srv.Name] = true
			servers = append(servers, srv)
		}
	}

	return servers
}

// parseMCPFile reads a single JSON file and extracts mcpServers entries.
func parseMCPFile(path string) []model.MCPServer {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var raw mcpFileFormat
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil
	}

	var servers []model.MCPServer
	for name, srv := range raw.MCPServers {
		servers = append(servers, model.MCPServer{
			Name:    name,
			Type:    srv.Type,
			URL:     srv.URL,
			Command: srv.Command,
			Args:    srv.Args,
		})
	}
	return servers
}

// parseClaudeJSON reads ~/.claude.json and extracts mcpServers for the specific project root.
func parseClaudeJSON(projectRoot string) []model.MCPServer {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	path := filepath.Join(home, ".claude.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var raw claudeJSONFormat
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil
	}

	var servers []model.MCPServer
	if p, ok := raw.Projects[projectRoot]; ok {
		for name, srv := range p.MCPServers {
			servers = append(servers, model.MCPServer{
				Name:    name,
				Type:    srv.Type,
				URL:     srv.URL,
				Command: srv.Command,
				Args:    srv.Args,
			})
		}
	}
	return servers
}
