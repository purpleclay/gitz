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
	assert.Equal(t, entries[0].Message, "pass tests")
	assert.Empty(t, entries[0].Tags)
	assert.ElementsMatch(t, []string{"HEAD -> new-feature", "origin/new-feature"}, entries[0].Branches)
	assert.False(t, entries[0].IsTrunk)
	assert.Equal(t, "new-feature", entries[0].HeadPointerRef)

	assert.Equal(t, entries[1].Message, "write tests for new feature")
	assert.Empty(t, entries[1].Tags)
	assert.Empty(t, entries[1].Branches)
	assert.False(t, entries[1].IsTrunk)
	assert.Empty(t, entries[1].HeadPointerRef)

	assert.Equal(t, entries[2].Message, "feat: improve existing cli documentation")
	assert.ElementsMatch(t, []string{"0.2.0", "v1"}, entries[2].Tags)
	assert.ElementsMatch(t, []string{"main", "origin/main"}, entries[2].Branches)
	assert.True(t, entries[2].IsTrunk)
	assert.Empty(t, entries[2].HeadPointerRef)

	assert.Equal(t, entries[3].Message, "docs: create initial mkdocs material documentation")
	assert.Empty(t, entries[3].Branches)
	assert.False(t, entries[3].IsTrunk)
	assert.Empty(t, entries[3].HeadPointerRef)

	assert.Equal(t, entries[4].Message, "feat: add secondary cli command to support filtering of results")
	assert.ElementsMatch(t, []string{"0.1.0"}, entries[4].Tags)
	assert.Empty(t, entries[4].Branches)
	assert.False(t, entries[4].IsTrunk)
	assert.Empty(t, entries[4].HeadPointerRef)

	assert.Equal(t, entries[5].Message, "feat: scaffold initial cli and add first command")
	assert.Empty(t, entries[5].Branches)
	assert.False(t, entries[5].IsTrunk)
	assert.Empty(t, entries[5].HeadPointerRef)
}

func TestParseLogMultiLineMode(t *testing.T) {
	log := `> (tag: 0.1.0, main, origin/main) fix: ensure parsing of multi-line commits is supported
> feat(deps): bump github.com/stretchr/testify from 1.8.1 to 1.8.2

Signed-off-by: dependabot[bot] <support@github.com>
Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`

	entries := gittest.ParseLog(log)

	require.Len(t, entries, 2)
	assert.Equal(t, "fix: ensure parsing of multi-line commits is supported", entries[0].Message)
	assert.ElementsMatch(t, []string{"0.1.0"}, entries[0].Tags)
	assert.Equal(t, `feat(deps): bump github.com/stretchr/testify from 1.8.1 to 1.8.2

Signed-off-by: dependabot[bot] <support@github.com>
Co-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>`, entries[1].Message)
}

func TestParseLogWithOptionalLeadingHash(t *testing.T) {
	log := `> b0d5429b967b9af0a0805fc2981b4420e10be38d feat: ensure parsing of optional leading hash is supported
> 58d708cb071df97e2561903aadcd4129419e9631 feat: include additional flag in pretty statements`

	entries := gittest.ParseLog(log)

	require.Len(t, entries, 2)
	assert.Equal(t, "b0d5429b967b9af0a0805fc2981b4420e10be38d", entries[0].Hash)
	assert.Equal(t, "b0d5429", entries[0].AbbrevHash)
	assert.Equal(t, "58d708cb071df97e2561903aadcd4129419e9631", entries[1].Hash)
	assert.Equal(t, "58d708c", entries[1].AbbrevHash)
}

func TestParseLogEmpty(t *testing.T) {
	entries := gittest.ParseLog("")
	assert.Empty(t, entries)
}

func TestParseLogTrimsSpaces(t *testing.T) {
	log := "   feat: testing if leading and trailing spaces are removed   "

	entries := gittest.ParseLog(log)

	require.Len(t, entries, 1)
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
				Message: "(tag: 0.1.0, HEAD -> main, main feat: this is a brand new feature",
			},
		},
		{
			name: "NoOpeningParentheses",
			log:  "HEAD -> main, main) ci: updated existing github workflow",
			expected: gittest.LogEntry{
				Message: "HEAD -> main, main) ci: updated existing github workflow",
			},
		},
		{
			name: "MismatchedParentheses",
			log:  "(tag: 0.2.0, HEAD -> main, main, new-feature)docs: include tests (and regressions) in guide",
			expected: gittest.LogEntry{
				Message:  "in guide",
				Tags:     []string{"0.2.0"},
				Branches: []string{"HEAD -> main", "main", "new-feature)docs: include tests (and regressions"},
			},
		},
		{
			name: "EmptyRefNames",
			log:  "() chore: add new issue template",
			expected: gittest.LogEntry{
				Message: "chore: add new issue template",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries := gittest.ParseLog(tt.log)

			require.Len(t, entries, 1)
			require.Equal(t, tt.expected.Message, entries[0].Message, "commit does not match")
			require.ElementsMatch(t, tt.expected.Tags, entries[0].Tags, "tags slice mismatch")
			require.ElementsMatch(t, tt.expected.Branches, entries[0].Branches, "branches slice mismatch")
		})
	}
}
