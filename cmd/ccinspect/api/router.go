package api

import (
	"encoding/json"
	"net/http"
)

// Router returns an HTTP handler with all /api/* routes registered.
func Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/scan", handleScan)
	mux.HandleFunc("GET /api/list", handleList)
	mux.HandleFunc("GET /api/global", handleGlobal)
	mux.HandleFunc("GET /api/diff", handleDiff)
	mux.HandleFunc("POST /api/copy", handleCopy)
	mux.HandleFunc("POST /api/clean", handleClean)
	mux.HandleFunc("POST /api/clean/teams", handleCleanTeams)
	mux.HandleFunc("GET /api/config", handleGetConfig)
	mux.HandleFunc("POST /api/config", handleSaveConfig)
	mux.HandleFunc("GET /api/fs/browse", handleFSBrowse)
	mux.HandleFunc("GET /api/fs/home", handleGetHomeDir)

	return mux
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}
