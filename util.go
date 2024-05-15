package git

import (
	"fmt"
	"strings"
)

func trim(strs ...string) []string {
	out := make([]string, 0, len(strs))
	for _, s := range strs {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			continue
		}

		out = append(out, trimmed)
	}

	return out
}

func trimAndPrefix(prefix string, strs ...string) []string {
	out := make([]string, 0, len(strs))
	for _, s := range strs {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			continue
		}

		if !strings.HasPrefix(trimmed, prefix) {
			trimmed = fmt.Sprintf("%s%s", prefix, trimmed)
		}
		out = append(out, trimmed)
	}

	return out
}

func reverse(strs ...string) []string {
	out := make([]string, 0, len(strs))
	for i := len(strs) - 1; i >= 0; i-- {
		out = append(out, strs[i])
	}

	return out
}
