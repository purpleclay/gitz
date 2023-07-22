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
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommit(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles("test.txt"))

	client, _ := git.NewClient()
	_, err := client.Commit("this is an example commit message")

	require.NoError(t, err)

	lastCommit := gittest.LastCommit(t)
	assert.Equal(t, gittest.DefaultAuthorName, lastCommit.AuthorName)
	assert.Equal(t, gittest.DefaultAuthorEmail, lastCommit.AuthorEmail)
	assert.Equal(t, "this is an example commit message", lastCommit.Message)
}

func TestCommitWithAllowEmpty(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	_, err := client.Commit("this will be a an empty commit", git.WithAllowEmpty())

	require.NoError(t, err)
}

func TestCommitWithNoGpgSign(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("test.txt"))
	gittest.ConfigSet(t, "user.signingkey", "DOES-NOT-EXIST", "commit.gpgsign", "true")
	gittest.StageFile(t, "test.txt")

	client, _ := git.NewClient()
	_, err := client.Commit("this will be a regular commit", git.WithNoGpgSign())

	require.NoError(t, err)
}

func TestCommitWithCommitConfig(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles("test.txt"))

	client, _ := git.NewClient()
	_, err := client.Commit("commit with inline options",
		git.WithCommitConfig("user.name", "bane", "user.email", "bane@dc.com"))

	require.NoError(t, err)
	lastCommit := gittest.LastCommit(t)
	assert.Equal(t, "bane", lastCommit.AuthorName)
	assert.Equal(t, "bane@dc.com", lastCommit.AuthorEmail)
}
