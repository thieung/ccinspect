package api

import (
	"net/http"

	"github.com/thieung/ccinspect/internal/config"
	"github.com/thieung/ccinspect/internal/parser"
	"github.com/thieung/ccinspect/internal/scanner"
)

func handleScan(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	cfg := config.Load()
	globalPath, _ := scanner.FindGlobal()
	claudePaths := scanner.FindClaudeDirs(cfg.ScanPaths, cfg.MaxDepth, cfg.ExcludePaths)
	inv := parser.BuildInventory(globalPath, claudePaths)

	var global *GlobalSummary
	if inv.Global != nil {
		g := inv.Global
		global = &GlobalSummary{
			Path:        g.Path,
			Skills:      len(g.Skills),
			Hooks:       len(g.Hooks),
			Agents:      len(g.Agents),
			Commands:    len(g.Commands),
			Rules:       len(g.Rules),
			Teams:       len(g.Teams),
			HasClaudeMD: g.HasClaudeMD,
		}
	}

	projects := make([]ProjectSummary, 0, len(inv.Projects))
	for _, p := range inv.Projects {
		projects = append(projects, ProjectSummary{
			Path:        p.Path,
			Skills:      len(p.Skills),
			Hooks:       len(p.Hooks),
			Agents:      len(p.Agents),
			Commands:    len(p.Commands),
			MCPServers:  len(p.MCPServers),
			HasClaudeMD: p.HasClaudeMD,
		})
	}

	writeJSON(w, http.StatusOK, ScanResponse{
		Projects: projects,
		Global:   global,
	})
}
