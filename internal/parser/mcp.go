package parser

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/thieunv/ccinspect/internal/model"
)

// ParseMCP reads .mcp.json from the project root (parent of .claude/).
func ParseMCP(claudePath string) []model.MCPServer {
	// .mcp.json is at project root, which is parent of .claude/
	projectRoot := filepath.Dir(claudePath)
	p := filepath.Join(projectRoot, ".mcp.json")

	data, err := os.ReadFile(p)
	if err != nil {
		return nil
	}

	var raw struct {
		MCPServers map[string]struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"mcpServers"`
	}
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
