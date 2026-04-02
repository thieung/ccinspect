package api

import (
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

type FSBrowseResponse struct {
	Entries []FSEntry `json:"entries"`
	Cwd     string    `json:"cwd"`
}

type FSEntry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size,omitempty"`
	ModTime int64  `json:"mod_time,omitempty"`
}

func handleFSBrowse(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "GET" {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		// Default to home directory
		home, err := os.UserHomeDir()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Cannot determine home directory")
			return
		}
		path = home
	}

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Cannot expand home directory")
			return
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	// Security: prevent accessing root-level system directories
	if !isPathAllowed(path) {
		writeError(w, http.StatusForbidden, "Access denied to this path")
		return
	}

	// Clean and resolve path
	path, err := filepath.Abs(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Invalid path")
		return
	}

	// Read directory entries
	entries, err := os.ReadDir(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Cannot read directory: "+err.Error())
		return
	}

	response := FSBrowseResponse{
		Cwd: path,
		Entries: make([]FSEntry, 0),
	}

	// Add parent directory entry if not at root
	parent := filepath.Dir(path)
	if parent != path {
		response.Entries = append(response.Entries, FSEntry{
			Name:  "..",
			Path:  parent,
			IsDir: true,
		})
	}

	// Add directory entries (directories first, then files)
	dirs := make([]FSEntry, 0)
	files := make([]FSEntry, 0)

	for _, entry := range entries {
		// Skip hidden files and symlinks
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())

		fsEntry := FSEntry{
			Name:    entry.Name(),
			Path:    fullPath,
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		}

		if entry.IsDir() {
			dirs = append(dirs, fsEntry)
		} else {
			files = append(files, fsEntry)
		}
	}

	response.Entries = append(response.Entries, dirs...)
	response.Entries = append(response.Entries, files...)

	writeJSON(w, http.StatusOK, response)
}

// isPathAllowed checks if the path is allowed to be accessed
// Prevents accessing system root and sensitive directories
func isPathAllowed(path string) bool {
	// Normalize path
	path = filepath.Clean(path)

	// Get OS-specific restrictions
	switch runtime.GOOS {
	case "windows":
		// Allow common drives (C:, D:, etc.) but not system roots
		if len(path) >= 2 && path[1] == ':' {
			// Allow drive roots, but not deep system paths
			denyList := []string{
				"\\windows\\system32",
				"\\program files\\windowsapps",
			}
			lowerPath := strings.ToLower(path)
			for _, deny := range denyList {
				if strings.Contains(lowerPath, deny) {
					return false
				}
			}
			return true
		}
	case "darwin", "linux":
		// Deny access to system root and sensitive paths
		denyList := []string{
			"/",
			"/etc",
			"/var",
			"/proc",
			"/sys",
			"/root",
			"/Library",
			"/System",
			"/Applications",
		}
		for _, deny := range denyList {
			if path == deny || strings.HasPrefix(path, deny+"/") {
				return false
			}
		}
	}

	return true
}

// handleGetHomeDir returns the user's home directory path
func handleGetHomeDir(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != "GET" {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Cannot determine home directory")
		return
	}

	// Get current user info
	currentUser, err := user.Current()
	username := "user"
	if err == nil {
		username = currentUser.Username
	}

	response := map[string]string{
		"home":     home,
		"username": username,
		"os":       runtime.GOOS,
	}

	writeJSON(w, http.StatusOK, response)
}
