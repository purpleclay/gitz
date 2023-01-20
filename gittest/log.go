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

package gittest

import (
	"bufio"
	"strings"
)

// LogEntry ...
type LogEntry struct {
	Commit string
	Tag    string
}

// ParseLog ...
func ParseLog(log string) []LogEntry {
	entries := make([]LogEntry, 0)

	scanner := bufio.NewScanner(strings.NewReader(log))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		entry := LogEntry{}

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "(tag:") {
			// Parse the tag from the log line and add it to the log entry
			tag, commit, _ := strings.Cut(line, ") ")
			entry.Commit = commit

			// Process the tag, and strip off the 'tag: ' prefix
			entry.Tag = strings.TrimLeft(tag, "(tag: ")
		} else {
			// Use the raw line of the log, as this will be a plain commit message
			entry.Commit = line
		}

		entries = append(entries, entry)
	}

	return entries
}
