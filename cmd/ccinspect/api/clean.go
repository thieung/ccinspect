package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func handleClean(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var req CleanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if req.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	target := expandCleanPath(req.Path)
	claudeDir := filepath.Join(target, ".claude")

	info, err := os.Stat(claudeDir)
	if err != nil || !info.IsDir() {
		writeError(w, http.StatusNotFound, "no .claude/ directory found at "+req.Path)
		return
	}

	var files []string
	filepath.Walk(claudeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if req.DryRun {
		writeJSON(w, http.StatusOK, CleanResponse{
			Status:     "ok",
			Message:    "dry run — no files removed",
			FilesCount: len(files),
			Files:      files,
			DryRun:     true,
		})
		return
	}

	// Remove .claude/
	if err := os.RemoveAll(claudeDir); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove .claude/: "+err.Error())
		return
	}

	// Remove .mcp.json if present
	mcpFile := filepath.Join(target, ".mcp.json")
	if _, err := os.Stat(mcpFile); err == nil {
		os.Remove(mcpFile)
	}

	writeJSON(w, http.StatusOK, CleanResponse{
		Status:     "ok",
		Message:    "cleaned .claude/ from " + req.Path,
		FilesCount: len(files),
		DryRun:     false,
	})
}

func handleCleanTeams(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var req CleanTeamsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "cannot determine home directory")
		return
	}
	teamsDir := filepath.Join(home, ".claude", "teams")

	entries, err := os.ReadDir(teamsDir)
	if err != nil {
		writeError(w, http.StatusNotFound, "no teams directory found")
		return
	}

	var toRemove []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		teamPath := filepath.Join(teamsDir, e.Name())
		hasConfig := false
		if _, err := os.Stat(filepath.Join(teamPath, "config.json")); err == nil {
			hasConfig = true
		}

		// Filter by team names if provided
		if len(req.TeamNames) > 0 {
			found := false
			for _, n := range req.TeamNames {
				if n == e.Name() {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if req.All || !hasConfig {
			toRemove = append(toRemove, e.Name())
		}
	}

	if req.DryRun {
		writeJSON(w, http.StatusOK, CleanTeamsResponse{
			Status:     "ok",
			Message:    "dry run — no teams removed",
			TeamsCount: len(toRemove),
			TeamNames:  toRemove,
			DryRun:     true,
		})
		return
	}

	for _, name := range toRemove {
		teamPath := filepath.Join(teamsDir, name)
		os.RemoveAll(teamPath)
	}

	writeJSON(w, http.StatusOK, CleanTeamsResponse{
		Status:     "ok",
		Message:    fmt.Sprintf("removed %d teams", len(toRemove)),
		TeamsCount: len(toRemove),
		TeamNames:  toRemove,
		DryRun:     false,
	})
}

func expandCleanPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
