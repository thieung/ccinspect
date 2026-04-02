package api

import (
	"encoding/json"
	"net/http"

	"github.com/thieung/ccinspect/internal/config"
)

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	cfg := config.Load()
	writeJSON(w, http.StatusOK, ConfigResponse{
		Config: ConfigData{
			ScanPaths:     cfg.ScanPaths,
			ExcludePaths:  cfg.ExcludePaths,
			MaxDepth:      cfg.MaxDepth,
			DefaultOutput: cfg.DefaultOutput,
		},
	})
}

func handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var req SaveConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	cfg := &config.Config{
		ScanPaths:     req.ScanPaths,
		ExcludePaths:  req.ExcludePaths,
		MaxDepth:      req.MaxDepth,
		DefaultOutput: req.DefaultOutput,
	}

	if err := config.Save(cfg); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save config: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
