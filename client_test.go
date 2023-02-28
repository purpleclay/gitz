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

func TestNewClientGitFound(t *testing.T) {
	client, err := git.NewClient()

	require.NoError(t, err)
	expected := gittest.MustExec(t, "git --version")
	assert.Equal(t, expected, client.Version())
}

func TestNewClientGitMissingError(t *testing.T) {
	// Temporarily remove git from the PATH
	t.Setenv("PATH", "/fake")

	client, err := git.NewClient()

	require.ErrorAs(t, err, &git.ErrGitMissing{})
	assert.EqualError(t, err, "git is not installed under the PATH environment variable. PATH resolves to /fake")
	assert.Nil(t, client)
}

func TestRepository(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	repo, err := client.Repository()

	require.NoError(t, err)
	assert.False(t, repo.DetachedHead)
	assert.False(t, repo.ShallowClone)
	assert.Equal(t, gittest.DefaultBranch, repo.DefaultBranch)
}

func TestRepositoryDetectsShallowClone(t *testing.T) {
	gittest.InitRepository(t, gittest.WithCloneDepth(1))

	client, _ := git.NewClient()
	repo, err := client.Repository()

	require.NoError(t, err)
	assert.True(t, repo.ShallowClone)
}

func TestRepositoryDetectsDetachedHead(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLocalCommits("chore: checking this out will force a detached head"))

	hash := latestCommitHash(t)
	gittest.Checkout(t, hash)

	client, _ := git.NewClient()
	repo, err := client.Repository()

	require.NoError(t, err)
	assert.True(t, repo.DetachedHead)
}

func TestRepositoryNotWorkingDirectory(t *testing.T) {
	nonWorkingDirectory(t)

	client, _ := git.NewClient()
	_, err := client.Repository()

	require.EqualError(t, err, "current working directory is not a git repository")
}
