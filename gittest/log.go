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

// LogEntry defines a single log entry from the history
// of a git repository
type LogEntry struct {
	// Commit contains the commit message
	Commit string

	// Tag contains a valid tag reference to an associated
	// commit within a log entry
	Tag string
}

// ParseLog will attempt to parse a log extract for a given repository
// into a series of commits and associated tags. The log will be returned
// in the chronological order provided.
//
// The log is expected to be in the following format:
//
//	(tag: 0.1.0) feat: improve existing cli documentation
//	docs: create initial mkdocs material documentation
//	feat: add secondary cli command to support filtering of results
//	feat: scaffold initial cli and add first command
//
// This is the equivalent to the format produced using the git command:
//
//	git log --pretty='format:%d %s'
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
