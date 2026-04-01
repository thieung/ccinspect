package output

import (
	"encoding/json"
	"fmt"
)

// RenderJSON outputs data as indented JSON to stdout.
func RenderJSON(data any) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(b)
}
