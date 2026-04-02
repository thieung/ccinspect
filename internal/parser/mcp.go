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
		Command string   `json:"command"`
		Args    []string `json:"args"`
	} `json:"mcpServers"`
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
			Command: srv.Command,
			Args:    srv.Args,
		})
	}
	return servers
}
