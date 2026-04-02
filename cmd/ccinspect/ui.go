package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thieung/ccinspect/internal/config"
	"github.com/thieung/ccinspect/internal/parser"
	"github.com/thieung/ccinspect/internal/scanner"
	"github.com/thieung/ccinspect/web"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the CCInspect Web UI dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Starting CCInspect Dashboard on http://localhost%s\n", addr)

		mux := http.NewServeMux()

		// Static UI distribution
		distFS, err := fs.Sub(web.DistFS, "dist")
		if err != nil {
			fmt.Println("Warning: Could not extract dist subfolder, using fallback")
		} else {
			mux.Handle("/", http.FileServer(http.FS(distFS)))
		}

		// API Handlers
		mux.HandleFunc("/api/scan", handleApiScan)
		mux.HandleFunc("/api/list", handleApiList)
		mux.HandleFunc("/api/global", handleApiGlobal)
		mux.HandleFunc("/api/diff", handleApiDiff)
		mux.HandleFunc("/api/copy", handleApiCopy)
		mux.HandleFunc("/api/clean", handleApiClean)

		return http.ListenAndServe(addr, mux)
	},
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleApiScan(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	cfg := config.Load()
	globalPath, _ := scanner.FindGlobal()
	claudePaths := scanner.FindClaudeDirs(cfg.ScanPaths, cfg.MaxDepth, cfg.ExcludePaths)
	inv := parser.BuildInventory(globalPath, claudePaths)

	json.NewEncoder(w).Encode(inv)
}

func handleApiGlobal(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	globalPath, _ := scanner.FindGlobal()
	inv := parser.BuildInventory(globalPath, []string{})

	json.NewEncoder(w).Encode(inv.Global)
}

func handleApiList(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "application/json")

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

	type EntityItem struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Source      string `json:"source"`
	}
	var results []EntityItem

	if entityType == "skills" {
		if inv.Global != nil {
			for _, s := range inv.Global.Skills {
				if prefix == "" || strings.HasPrefix(s.Name, prefix) {
					results = append(results, EntityItem{Name: s.Name, Description: s.Description, Source: "global"})
				}
			}
		}
		for _, p := range inv.Projects {
			for _, s := range p.Skills {
				if prefix == "" || strings.HasPrefix(s.Name, prefix) {
					results = append(results, EntityItem{Name: s.Name, Description: s.Description, Source: p.Path})
				}
			}
		}
	} else if entityType == "hooks" {
		if inv.Global != nil {
			for _, h := range inv.Global.Hooks {
				results = append(results, EntityItem{Name: h.Event, Description: "Global Hook", Source: h.Command})
			}
		}
		for _, p := range inv.Projects {
			for _, h := range p.Hooks {
				results = append(results, EntityItem{Name: h.Event, Description: "Project Hook", Source: h.Command})
			}
		}
	}

	json.NewEncoder(w).Encode(results)
}

func handleApiDiff(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "not_implemented"}`))
}

func handleApiCopy(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Type  string   `json:"type"`
		From  string   `json:"from"`
		To    string   `json:"to"`
		Names []string `json:"names"`
		Force bool     `json:"force"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// returning ok for now
	w.Write([]byte(`{"status": "ok"}`))
}

func handleApiClean(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

func init() {
	uiCmd.Flags().IntP("port", "p", 8080, "Port to run the Web UI on")
	rootCmd.AddCommand(uiCmd)
}
