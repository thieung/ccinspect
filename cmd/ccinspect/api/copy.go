package api

import (
	"encoding/json"
	"net/http"

	"github.com/thieung/ccinspect/internal/copier"
)

func handleCopy(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var req CopyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if req.Type == "" || req.From == "" || req.To == "" {
		writeError(w, http.StatusBadRequest, "missing required fields: type, from, to")
		return
	}
	if len(req.Names) == 0 {
		req.Names = []string{"all"}
	}

	results, err := copier.CopyEntities(req.From, req.To, req.Type, req.Names, req.Force, req.DryRun)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	copyResults := make([]CopyResult, len(results))
	for i, r := range results {
		copyResults[i] = CopyResult{
			Name:   r.Name,
			From:   r.From,
			To:     r.To,
			Status: r.Status,
			Detail: r.Detail,
		}
	}

	writeJSON(w, http.StatusOK, CopyResponse{
		Status:  "ok",
		Results: copyResults,
	})
}
