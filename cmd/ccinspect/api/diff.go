package api

import (
	"net/http"

	"github.com/thieung/ccinspect/internal/config"
	"github.com/thieung/ccinspect/internal/parser"
	"github.com/thieung/ccinspect/internal/scanner"
)

func handleDiff(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	pathA := r.URL.Query().Get("a")
	pathB := r.URL.Query().Get("b")
	entityType := r.URL.Query().Get("type")
	if entityType == "" {
		entityType = "skills"
	}

	if pathA == "" || pathB == "" {
		writeError(w, http.StatusBadRequest, "missing required params: a and b")
		return
	}

	skillsA := resolveEntities(pathA, entityType)
	skillsB := resolveEntities(pathB, entityType)

	leftMap := make(map[string]bool)
	rightMap := make(map[string]bool)

	var leftName, rightName string
	if pathA == "global" {
		leftName = "global"
	} else {
		leftName = pathA
	}
	if pathB == "global" {
		rightName = "global"
	} else {
		rightName = pathB
	}

	for _, s := range skillsA {
		leftMap[s.Name] = true
	}
	for _, s := range skillsB {
		rightMap[s.Name] = true
	}

	var added, removed []DiffEntry
	for _, s := range skillsA {
		if !rightMap[s.Name] {
			removed = append(removed, DiffEntry{Name: s.Name, Description: s.Description, Side: "left"})
		}
	}
	for _, s := range skillsB {
		if !leftMap[s.Name] {
			added = append(added, DiffEntry{Name: s.Name, Description: s.Description, Side: "right"})
		}
	}

	writeJSON(w, http.StatusOK, DiffResponse{
		LeftName:   leftName,
		RightName:  rightName,
		EntityType: entityType,
		Added:      added,
		Removed:    removed,
		Changed:    nil,
	})
}

func resolveEntities(path, entityType string) []struct{ Name string; Description string } {
	var result []struct{ Name string; Description string }

	if path == "global" {
		globalPath, _ := scanner.FindGlobal()
		if globalPath == "" {
			return result
		}
		inv := parser.BuildInventory(globalPath, nil)
		if inv.Global == nil {
			return result
		}
		switch entityType {
		case "skills":
			for _, s := range inv.Global.Skills {
				result = append(result, struct{ Name string; Description string }{s.Name, s.Description})
			}
		case "hooks":
			for _, h := range inv.Global.Hooks {
				result = append(result, struct{ Name string; Description string }{h.Event, h.Command})
			}
		case "agents":
			for _, e := range inv.Global.Agents {
				result = append(result, struct{ Name string; Description string }{e.Name, ""})
			}
		case "commands":
			for _, e := range inv.Global.Commands {
				result = append(result, struct{ Name string; Description string }{e.Name, ""})
			}
		}
		return result
	}

	cfg := config.Load()
	dirs := scanner.FindClaudeDirs([]string{path}, 1, cfg.ExcludePaths)
	if len(dirs) == 0 {
		return result
	}
	inv := parser.BuildInventory("", dirs)
	if len(inv.Projects) == 0 {
		return result
	}
	p := inv.Projects[0]
	switch entityType {
	case "skills":
		for _, s := range p.Skills {
			result = append(result, struct{ Name string; Description string }{s.Name, s.Description})
		}
	case "hooks":
		for _, h := range p.Hooks {
			result = append(result, struct{ Name string; Description string }{h.Event, h.Command})
		}
	case "agents":
		for _, e := range p.Agents {
			result = append(result, struct{ Name string; Description string }{e.Name, ""})
		}
	case "commands":
		for _, e := range p.Commands {
			result = append(result, struct{ Name string; Description string }{e.Name, ""})
		}
	case "mcp":
		for _, m := range p.MCPServers {
			result = append(result, struct{ Name string; Description string }{m.Name, m.URL})
		}
	}
	return result
}
