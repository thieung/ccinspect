package api

import (
	"net/http"
	"strings"

	"github.com/thieung/ccinspect/internal/config"
	"github.com/thieung/ccinspect/internal/parser"
	"github.com/thieung/ccinspect/internal/scanner"
)

func handleList(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	entityType := r.URL.Query().Get("type")
	globalOnly := r.URL.Query().Get("global") == "true"
	prefix := r.URL.Query().Get("prefix")

	cfg := config.Load()
	globalPath, _ := scanner.FindGlobal()
	var claudePaths []string
	if !globalOnly {
		claudePaths = scanner.FindClaudeDirs(cfg.ScanPaths, cfg.MaxDepth, cfg.ExcludePaths)
	}

	inv := parser.BuildInventory(globalPath, claudePaths)
	var results []EntityItem

	switch entityType {
	case "skills":
		for _, s := range inv.Global.Skills {
			if prefix == "" || strings.HasPrefix(s.Name, prefix) {
				results = append(results, EntityItem{Name: s.Name, Description: s.Description, Source: "global", Type: "skill"})
			}
		}
		for _, p := range inv.Projects {
			for _, s := range p.Skills {
				if prefix == "" || strings.HasPrefix(s.Name, prefix) {
					results = append(results, EntityItem{Name: s.Name, Description: s.Description, Source: p.Path, Type: "skill"})
				}
			}
		}
	case "hooks":
		for _, h := range inv.Global.Hooks {
			results = append(results, EntityItem{Name: h.Event, Description: h.Command, Source: "global", Type: "hook"})
		}
		for _, p := range inv.Projects {
			for _, h := range p.Hooks {
				results = append(results, EntityItem{Name: h.Event, Description: h.Command, Source: p.Path, Type: "hook"})
			}
		}
	case "agents":
		for _, e := range inv.Global.Agents {
			results = append(results, EntityItem{Name: e.Name, Source: "global", Type: "agent"})
		}
		for _, p := range inv.Projects {
			for _, e := range p.Agents {
				results = append(results, EntityItem{Name: e.Name, Source: p.Path, Type: "agent"})
			}
		}
	case "commands":
		for _, e := range inv.Global.Commands {
			results = append(results, EntityItem{Name: e.Name, Source: "global", Type: "command"})
		}
		for _, p := range inv.Projects {
			for _, e := range p.Commands {
				results = append(results, EntityItem{Name: e.Name, Source: p.Path, Type: "command"})
			}
		}
	case "rules":
		for _, e := range inv.Global.Rules {
			results = append(results, EntityItem{Name: e.Name, Source: "global", Type: "rule"})
		}
	case "mcp":
		for _, p := range inv.Projects {
			for _, m := range p.MCPServers {
				results = append(results, EntityItem{Name: m.Name, Source: p.Path, Type: "mcp"})
			}
		}
	case "teams":
		for _, e := range inv.Global.Teams {
			results = append(results, EntityItem{Name: e.Name, Source: "global", Type: "team"})
		}
	default:
		writeError(w, http.StatusBadRequest, "invalid type: must be skills|hooks|agents|commands|rules|mcp|teams")
		return
	}

	writeJSON(w, http.StatusOK, ListResponse{Entities: results})
}

func handleGlobal(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	globalPath, _ := scanner.FindGlobal()
	inv := parser.BuildInventory(globalPath, nil)

	if inv.Global == nil {
		writeJSON(w, http.StatusOK, map[string]any{"global": nil})
		return
	}

	g := inv.Global
	writeJSON(w, http.StatusOK, map[string]any{
		"global": map[string]any{
			"path":        g.Path,
			"skills":      g.Skills,
			"hooks":       g.Hooks,
			"agents":      g.Agents,
			"commands":    g.Commands,
			"rules":       g.Rules,
			"teams":       g.Teams,
			"settings":    g.Settings,
			"has_claude_md": g.HasClaudeMD,
		},
	})
}
