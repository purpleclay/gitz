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

package gittest_test

import (
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLog(t *testing.T) {
	log := `(HEAD -> new-feature, origin/new-feature) pass tests
write tests for new feature
(tag: 0.2.0, tag: v1, main, origin/main) feat: improve existing cli documentation
docs: create initial mkdocs material documentation
(tag: 0.1.0) feat: add secondary cli command to support filtering of results
feat: scaffold initial cli and add first command`

	entries := gittest.ParseLog(log)

	require.Len(t, entries, 6)
	assert.Equal(t, entries[0].Commit, "pass tests")
	assert.Equal(t, entries[0].Message, "pass tests")
	assert.Empty(t, entries[0].Tag)
	assert.Empty(t, entries[0].Tags)
	assert.ElementsMatch(t, []string{"HEAD -> new-feature", "origin/new-feature"}, entries[0].Branches)
	assert.False(t, entries[0].IsTrunk)
	assert.Equal(t, "new-feature", entries[0].HeadPointerRef)

	assert.Equal(t, entries[1].Commit, "write tests for new feature")
	assert.Equal(t, entries[1].Message, "write tests for new feature")
	assert.Empty(t, entries[1].Tag)
	assert.Empty(t, entries[1].Tags)
	assert.Empty(t, entries[1].Branches)
	assert.False(t, entries[1].IsTrunk)
	assert.Empty(t, entries[1].HeadPointerRef)

	assert.Equal(t, entries[2].Commit, "feat: improve existing cli documentation")
	assert.Equal(t, entries[2].Message, "feat: improve existing cli documentation")
	assert.Equal(t, "0.2.0", entries[2].Tag)
	assert.ElementsMatch(t, []string{"0.2.0", "v1"}, entries[2].Tags)
	assert.ElementsMatch(t, []string{"main", "origin/main"}, entries[2].Branches)
	assert.True(t, entries[2].IsTrunk)
	assert.Empty(t, entries[2].HeadPointerRef)

	assert.Equal(t, entries[3].Commit, "docs: create initial mkdocs material documentation")
	assert.Equal(t, entries[3].Message, "docs: create initial mkdocs material documentation")
	assert.Empty(t, entries[3].Tag)
	assert.Empty(t, entries[3].Branches)
	assert.False(t, entries[3].IsTrunk)
	assert.Empty(t, entries[3].HeadPointerRef)

	assert.Equal(t, entries[4].Commit, "feat: add secondary cli command to support filtering of results")
	assert.Equal(t, entries[4].Message, "feat: add secondary cli command to support filtering of results")
	assert.Equal(t, "0.1.0", entries[4].Tag)
	assert.ElementsMatch(t, []string{"0.1.0"}, entries[4].Tags)
	assert.Empty(t, entries[4].Branches)
	assert.False(t, entries[4].IsTrunk)
	assert.Empty(t, entries[4].HeadPointerRef)

	assert.Equal(t, entries[5].Commit, "feat: scaffold initial cli and add first command")
	assert.Equal(t, entries[5].Message, "feat: scaffold initial cli and add first command")
	assert.Empty(t, entries[5].Tag)
	assert.Empty(t, entries[5].Branches)
	assert.False(t, entries[5].IsTrunk)
	assert.Empty(t, entries[5].HeadPointerRef)
}

func TestParseLogEmpty(t *testing.T) {
	entries := gittest.ParseLog("")
	assert.Empty(t, entries)
}

func TestParseLogTrimsSpaces(t *testing.T) {
	log := "   feat: testing if leading and trailing spaces are removed   "

	entries := gittest.ParseLog(log)

	require.Len(t, entries, 1)
	assert.Equal(t, entries[0].Commit, "feat: testing if leading and trailing spaces are removed")
	assert.Equal(t, entries[0].Message, "feat: testing if leading and trailing spaces are removed")
}

func TestParseLogMalformedLine(t *testing.T) {
	tests := []struct {
		name     string
		log      string
		expected gittest.LogEntry
	}{
		{
			name: "NoClosingParentheses",
			log:  "(tag: 0.1.0, HEAD -> main, main feat: this is a brand new feature",
			expected: gittest.LogEntry{
				Commit:  "(tag: 0.1.0, HEAD -> main, main feat: this is a brand new feature",
				Message: "(tag: 0.1.0, HEAD -> main, main feat: this is a brand new feature",
			},
		},
		{
			name: "NoOpeningParentheses",
			log:  "HEAD -> main, main) ci: updated existing github workflow",
			expected: gittest.LogEntry{
				Commit:  "HEAD -> main, main) ci: updated existing github workflow",
				Message: "HEAD -> main, main) ci: updated existing github workflow",
			},
		},
		{
			name: "MismatchedParentheses",
			log:  "(tag: 0.2.0, HEAD -> main, main, new-feature)docs: include tests (and regressions) in guide",
			expected: gittest.LogEntry{
				Commit:   "in guide",
				Message:  "in guide",
				Tag:      "0.2.0",
				Tags:     []string{"0.2.0"},
				Branches: []string{"HEAD -> main", "main", "new-feature)docs: include tests (and regressions"},
			},
		},
		{
			name: "EmptyRefNames",
			log:  "() chore: add new issue template",
			expected: gittest.LogEntry{
				Commit:  "chore: add new issue template",
				Message: "chore: add new issue template",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries := gittest.ParseLog(tt.log)

			require.Len(t, entries, 1)
			require.Equal(t, tt.expected.Commit, entries[0].Commit, "commit does not match")
			require.Equal(t, tt.expected.Message, entries[0].Message, "commit does not match")
			require.Equal(t, tt.expected.Tag, entries[0].Tag, "tag does not match")
			require.ElementsMatch(t, tt.expected.Tags, entries[0].Tags, "tags slice mismatch")
			require.ElementsMatch(t, tt.expected.Branches, entries[0].Branches, "branches slice mismatch")
		})
	}
}
