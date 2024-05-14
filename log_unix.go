//go:build !windows
// +build !windows

package git

import "strings"

func cleanLineEndings(log string) string {
	return strings.ReplaceAll(log, "\r\n", "\n")
}
