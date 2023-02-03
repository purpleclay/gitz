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
	"os"
	"path/filepath"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: this should be deleted

func tempFile(t *testing.T, path, content string) {
	t.Helper()

	err := os.MkdirAll(filepath.Dir(path), 0o755)
	require.NoError(t, err)

	err = os.WriteFile(path, []byte(content), 0o644)
	require.NoError(t, err)

	t.Cleanup(func() {
		// Check for the files existence before attempting to remove it. Depending
		// on cleanup order, it may have already been removed
		if _, err := os.Stat(path); err != nil {
			require.NoError(t, os.RemoveAll(path))
		}
	})
}

func TestCommit(t *testing.T) {
	gittest.InitRepository(t)
	// TODO: add a new option to create the repository with temporary files, this can also be staged
	// WithFiles
	// WithStagedFiles
	tempFile(t, "test.txt", "this is a test")
	gittest.StageFile(t, "test.txt")

	client := git.NewClient()
	err := client.Commit("this is an example commit message")

	require.NoError(t, err)

	out := gittest.LastCommit(t)
	assert.Contains(t, out, gittest.DefaultAuthorLog)
	assert.Contains(t, out, "this is an example commit message")
}

func TestCommitCleanWorkingTree(t *testing.T) {
	gittest.InitRepository(t)

	client := git.NewClient()
	err := client.Commit("this is an example commit message")

	require.ErrorContains(t, err, "nothing to commit, working tree clean")
}
