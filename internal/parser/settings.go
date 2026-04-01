package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ParseSettings reads a settings JSON file and returns raw map.
func ParseSettings(claudePath string, filename string) (map[string]any, error) {
	p := filepath.Join(claudePath, filename)
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}
