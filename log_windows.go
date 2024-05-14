//go:build windows
// +build windows

package git

func cleanLineEndings(log string) string {
	// Mixed line endings don't appear within Windows
	return log
}
