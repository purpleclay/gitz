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

package git_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	log := `fix: parsing error when input string is too long
ci: extend the existing build workflow to include integration tests
docs: create initial mkdocs material documentation
feat: add second operation to library
feat: add first operation to library`

	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	out, err := client.Log()
	require.NoError(t, err)

	lines := countLogLines(t, out.Raw)
	require.Equal(t, 6, lines)
	require.Equal(t, 6, len(out.Commits))

	assert.Contains(t, out.Raw, "fix: parsing error when input string is too long")
	assert.Equal(t, out.Commits[0].Message, "fix: parsing error when input string is too long")

	assert.Contains(t, out.Raw, "ci: extend the existing build workflow to include integration tests")
	assert.Equal(t, out.Commits[1].Message, "ci: extend the existing build workflow to include integration tests")

	assert.Contains(t, out.Raw, "docs: create initial mkdocs material documentation")
	assert.Equal(t, out.Commits[2].Message, "docs: create initial mkdocs material documentation")

	assert.Contains(t, out.Raw, "feat: add second operation to library")
	assert.Equal(t, out.Commits[3].Message, "feat: add second operation to library")

	assert.Contains(t, out.Raw, "feat: add first operation to library")
	assert.Equal(t, out.Commits[4].Message, "feat: add first operation to library")

	assert.Contains(t, out.Raw, gittest.InitialCommit)
	assert.Equal(t, out.Commits[5].Message, gittest.InitialCommit)
}

// A utility function that will scan the raw output from a git log and
// count all of the returned log lines. It is important to note, that
// in some scenarios the log will contain the [gittest.InitialCommit]
// used to initialize the repository
func countLogLines(t *testing.T, log string) int {
	scanner := bufio.NewScanner(strings.NewReader(log))
	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}

// TODO: test function is parsing log line correctly. Use initial repository and verify the hash

// TestLogWithRawOnly

func TestLogWithRef(t *testing.T) {
	log := `(tag: 0.1.1) fix: unexpected bytes in message while parsing
(tag: 0.1.0) docs: create initial mkdocs material documentation
feat: build exciting new library`

	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	out, err := client.Log(git.WithRef("0.1.0"))
	require.NoError(t, err)

	lines := countLogLines(t, out.Raw)
	require.Equal(t, 3, lines)

	assert.Contains(t, out.Raw, "docs: create initial mkdocs material documentation")
	assert.Contains(t, out.Raw, "feat: build exciting new library")
	assert.Contains(t, out.Raw, gittest.InitialCommit)
}

func TestLogWithRefRange(t *testing.T) {
	log := `(tag: 0.2.0) feat: add ability to filter on results
(tag: 0.1.1) fix: unexpected bytes in message while parsing
docs: update documentation to include fix
(tag: 0.1.0) docs: create initial mkdocs material documentation
feat: build exciting new library`

	tests := []struct {
		name            string
		fromRef         string
		toRef           string
		expectedLines   int
		expectedCommits []string
	}{
		{
			name:          "FromAndToRefsProvided",
			fromRef:       "0.1.1",
			toRef:         "0.1.0",
			expectedLines: 2,
			expectedCommits: []string{
				"fix: unexpected bytes in message while parsing",
				"docs: update documentation to include fix",
			},
		},
		{
			name:          "FromRefOnly",
			fromRef:       "0.1.0",
			expectedLines: 3,
			expectedCommits: []string{
				"docs: create initial mkdocs material documentation",
				"feat: build exciting new library",
				gittest.InitialCommit,
			},
		},
		{
			name:          "ToRefOnly",
			toRef:         "0.1.1",
			expectedLines: 1,
			expectedCommits: []string{
				"feat: add ability to filter on results",
			},
		},
		{
			name:          "TrimsWhitespaceAroundRefs",
			fromRef:       "  0.2.0  ",
			toRef:         "  0.1.1  ",
			expectedLines: 1,
			expectedCommits: []string{
				"feat: add ability to filter on results",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gittest.InitRepository(t, gittest.WithLog(log))

			client, _ := git.NewClient()
			out, err := client.Log(git.WithRefRange(tt.fromRef, tt.toRef))
			require.NoError(t, err)

			lines := countLogLines(t, out.Raw)
			require.Equal(t, tt.expectedLines, lines)

			for _, commit := range tt.expectedCommits {
				require.Contains(t, out.Raw, commit)
			}
		})
	}
}

func TestLogWithPaths(t *testing.T) {
	gittest.InitRepository(t,
		gittest.WithLocalCommits("this should not appear in the log"),
		gittest.WithStagedFiles("dir1/a.txt", "dir2/b.txt"))

	gittest.Commit(t, "feat: include both dir1/a.txt and dir2/b.txt")
	overwriteFile(t, "dir1/a.txt", "Help, I have been overwritten!")
	gittest.StageFile(t, "dir1/a.txt")
	gittest.Commit(t, "fix: changed file dir1/a.txt")

	client, _ := git.NewClient()
	out, err := client.Log(git.WithPaths("dir1"))
	require.NoError(t, err)

	lines := countLogLines(t, out.Raw)
	require.Equal(t, 2, lines)
	assert.Contains(t, out.Raw, "fix: changed file dir1/a.txt")
	assert.Contains(t, out.Raw, "feat: include both dir1/a.txt and dir2/b.txt")
}

func overwriteFile(t *testing.T, path, content string) {
	t.Helper()

	fi, err := os.Create(path)
	require.NoError(t, err)
	defer fi.Close()

	fi.WriteString(content)
}
