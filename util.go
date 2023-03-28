/*
Copyright (c) 2023 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package git

import (
	"fmt"
	"strings"
)

// Trim iterates through a slice, trimming leading and trailing
// whitespace from each string. Empty strings are ignored
// and removed from the slice
func Trim(strs ...string) []string {
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

// TrimAndPrefix iterates through a slice, trimming leading and
// trailing whitespace from each string before appending the
// provided prefix. Empty strings are ignored and removed from
// the slice
func TrimAndPrefix(prefix string, strs ...string) []string {
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

// TrimAndRemove iterates through a slice, trimming leading and
// trailing whitespace from each string. Strings that are empty
// or match the removal string, are removed from the slice
func TrimAndRemove(rem string, strs ...string) []string {
	out := make([]string, 0, len(strs))
	for _, s := range strs {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" || trimmed == rem {
			continue
		}

		out = append(out, trimmed)
	}

	return out
}
