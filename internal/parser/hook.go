package parser

import (
	"fmt"

	"github.com/thieung/ccinspect/internal/model"
)

// ParseHooks extracts hook definitions from a parsed settings map.
func ParseHooks(settings map[string]any) []model.Hook {
	if settings == nil {
		return nil
	}
	hooksRaw, ok := settings["hooks"]
	if !ok {
		return nil
	}
	hooksMap, ok := hooksRaw.(map[string]any)
	if !ok {
		return nil
	}

	var hooks []model.Hook
	for event, matchersRaw := range hooksMap {
		matchers, ok := matchersRaw.([]any)
		if !ok {
			continue
		}
		for _, matcherRaw := range matchers {
			matcher, ok := matcherRaw.(map[string]any)
			if !ok {
				continue
			}
			matcherName, _ := matcher["matcher"].(string)
			hooksArr, ok := matcher["hooks"].([]any)
			if !ok {
				continue
			}
			for _, hookRaw := range hooksArr {
				hookMap, ok := hookRaw.(map[string]any)
				if !ok {
					continue
				}
				cmd, _ := hookMap["command"].(string)
				typ, _ := hookMap["type"].(string)
				if typ == "" {
					typ = "command"
				}
				hooks = append(hooks, model.Hook{
					Event:   event,
					Matcher: matcherName,
					Command: truncate(cmd, 60),
					Type:    typ,
				})
			}
		}
	}
	return hooks
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return fmt.Sprintf("%s...", s[:max-3])
}
