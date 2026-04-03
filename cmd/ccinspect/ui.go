package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/thieung/ccinspect/cmd/ccinspect/api"
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

		// API routes
		mux.Handle("/api/", api.Router())

		return http.ListenAndServe(addr, mux)
	},
}

func init() {
	uiCmd.Flags().IntP("port", "p", 8080, "Port to run the Web UI on")
	rootCmd.AddCommand(uiCmd)
}
