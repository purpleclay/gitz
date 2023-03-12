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

// LogEntry represents a single log entry from the history
// of a git repository
type LogEntry struct {
	// Commit contains the commit message
	Commit string

	// Tag contains a valid tag reference to an associated
	// commit within a log entry. If using multiple tags,
	// only the first will be referenced
	//
	// Deprecated: Use [gittest.LogEntry.Tags] instead.
	Tag string

	// Tags contains all tag references that are associated
	// with the current commit
	Tags []string

	// Branches contains the name of all branches (local and remote)
	// that are associated with the current commit
	Branches []string
}

// ParseLog will attempt to parse a log extract from a given repository
// into a series of commits, branches and tags. The log will be returned
// in the chronological order provided. The parser is designed to not
// error and parses each line with best endeavors.
//
// The log is expected to be in the following format:
//
//	(HEAD -> new-feature, origin/new-feature) pass tests
//	write tests for new feature
//	(tag: 0.2.0, tag: v1, main, origin/main) feat: improve existing cli documentation
//	docs: create initial mkdocs material documentation
//	(tag: 0.1.0) feat: add secondary cli command to support filtering of results
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
		line := scanner.Text()
		line = strings.TrimSpace(line)

		entry := LogEntry{Commit: line}
		if strings.HasPrefix(line, "(") {
			// Cut based on the first occurrence of a closing parentheses, if one doesn't
			// exist, then append the line as a raw log entry
			refNames, msg, found := strings.Cut(line, ") ")
			if !found {
				goto append
			}

			entry.Commit = msg

			// Process the comma separated list of ref names, preceding the commit message
			for _, ref := range strings.Split(refNames[1:], ",") {
				cleanedRef := strings.TrimSpace(ref)

				if cleanedRef == "" {
					continue
				}

				if strings.HasPrefix(cleanedRef, "tag: ") {
					entry.Tags = append(entry.Tags, strings.TrimPrefix(cleanedRef, "tag: "))
				} else {
					entry.Branches = append(entry.Branches, cleanedRef)
				}
			}

			// For backwards compatibility, store a reference to the first tag
			if len(entry.Tags) > 0 {
				entry.Tag = entry.Tags[0]
			}
		}

	append:
		entries = append(entries, entry)
	}

	return entries
}
