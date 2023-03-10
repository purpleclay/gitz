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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Creates the expected git status message of an untracked file, as based
// on the official git documentation: https://git-scm.com/docs/git-status#_short_format
//
//	?? file
func statusUntracked(file string) string {
	return fmt.Sprintf("?? %s", file)
}

// Creates the expected git status message of a staged file, as based
// on the official git documentation: https://git-scm.com/docs/git-status#_short_format
//
//	A  file
func statusAdded(file string) string {
	return fmt.Sprintf("A  %s", file)
}

func TestInitRepositoryConfigSet(t *testing.T) {
	gittest.InitRepository(t)

	out, err := exec.Command("git", "config", "--list").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), fmt.Sprintf("user.name=%s", gittest.DefaultAuthorName))
	assert.Contains(t, string(out), fmt.Sprintf("user.email=%s", gittest.DefaultAuthorEmail))
}

func TestInitRepositoryDefaultBranchSet(t *testing.T) {
	gittest.InitRepository(t)

	out, err := exec.Command("git", "branch").CombinedOutput()
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("* %s\n", gittest.DefaultBranch), string(out))
}

func TestInitRepositoryWithLog(t *testing.T) {
	log := `(tag: 0.1.0) feat: this is a brand new feature
ci: include github workflow`
	gittest.InitRepository(t, gittest.WithLog(log))

	out, err := exec.Command("git", "log", "-n2", "--oneline").CombinedOutput()
	require.NoError(t, err)

	tag, err := exec.Command("git", "tag").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), "feat: this is a brand new feature")
	assert.Contains(t, string(out), "ci: include github workflow")
	assert.Contains(t, string(tag), "0.1.0")
}

func TestInitRepositoryWithFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("a.txt", "b.txt"))

	out, err := exec.Command("git", "status", "--porcelain").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), statusUntracked("a.txt"))
	assert.Contains(t, string(out), statusUntracked("b.txt"))
}

func TestInitRepositoryWithStagedFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles("c.txt", "dir/d.txt"))

	out, err := exec.Command("git", "status", "--porcelain").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), statusAdded("c.txt"))
	assert.Contains(t, string(out), statusAdded("dir/d.txt"))
}

func TestInitRepositoryWithLocalCommits(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLocalCommits("local commit 1", "local commit 2"))

	out, err := exec.Command("git", "log", "--oneline").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), "local commit 1")
	assert.Contains(t, string(out), "local commit 2")

	out, err = exec.Command("git", "log", "--oneline", gittest.DefaultRemoteBranch).CombinedOutput()
	require.NoError(t, err)

	assert.NotContains(t, string(out), "local commit 1")
	assert.NotContains(t, string(out), "local commit 2")
}

func TestWithRemoteLog(t *testing.T) {
	log := "this is a remote commit"
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	localLog, err := exec.Command("git", "log", "-n1", "--oneline").CombinedOutput()
	require.NoError(t, err)
	assert.NotContains(t, string(localLog), log)

	_, err = exec.Command("git", "pull").CombinedOutput()
	require.NoError(t, err)

	localLog, err = exec.Command("git", "log", "-n1", "--oneline").CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(localLog), log)
}

func TestWithCloneDepth(t *testing.T) {
	log := `feat: this is commit number 3
feat: this is commit number 2
feat: this is commit number 1`

	gittest.InitRepository(t, gittest.WithLog(log), gittest.WithCloneDepth(1))

	localLog, err := exec.Command("git", "log", "-n4", "--oneline").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(localLog), "feat: this is commit number 3")
	assert.NotContains(t, string(localLog), "feat: this is commit number 2")
	assert.NotContains(t, string(localLog), "feat: this is commit number 1")
	assert.NotContains(t, string(localLog), gittest.InitialCommit)
}

func TestExecHasRawGitOutput(t *testing.T) {
	out, err := gittest.Exec(t, "git --version")

	require.NoError(t, err)
	assert.Contains(t, out, "git version")
}

func TestExecReturnsClientError(t *testing.T) {
	_, err := gittest.Exec(t, "git unknown")

	require.ErrorContains(t, err, "git: 'unknown' is not a git command")
}

func TestMustExecHasRawGitOutput(t *testing.T) {
	out := gittest.MustExec(t, "git --version")

	assert.Contains(t, out, "git version")
}

func TestTags(t *testing.T) {
	gittest.InitRepository(t)

	_, err := exec.Command("git", "tag", "0.1.0").CombinedOutput()
	require.NoError(t, err)

	out := gittest.Tags(t)
	assert.Contains(t, out, "refs/tags/0.1.0")
}

func TestRemoteTags(t *testing.T) {
	gittest.InitRepository(t)

	_, err := exec.Command("git", "tag", "0.2.0").CombinedOutput()
	require.NoError(t, err)

	_, err = exec.Command("git", "push", "origin", "0.2.0").CombinedOutput()
	require.NoError(t, err)

	out := gittest.RemoteTags(t)
	assert.Contains(t, out, "refs/tags/0.2.0")
}

func TestStageFile(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("test.txt"))

	gittest.StageFile(t, "test.txt")

	out, err := exec.Command("git", "status", "--porcelain").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), "A  test.txt")
}

func TestCommit(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles("file.txt"))

	gittest.Commit(t, "include file.txt")

	out, err := exec.Command("git", "log", "-n1", "--oneline").CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(out), "include file.txt")
}

func TestLastCommit(t *testing.T) {
	gittest.InitRepository(t)

	_, err := exec.Command("git", "commit", "--allow-empty", "-m", "this is a test").CombinedOutput()
	require.NoError(t, err)

	log := gittest.LastCommit(t)
	assert.Contains(t, log, "this is a test")
}

func TestPorcelainStatus(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("file1.txt", "file2.txt"))

	status := gittest.PorcelainStatus(t)
	assert.Equal(t, "?? file1.txt\n?? file2.txt", status)
}

func TestLogRemote(t *testing.T) {
	gittest.InitRepository(t)

	_, err := exec.Command("git", "commit", "--allow-empty", "-m", "this commit is on the remote").CombinedOutput()
	require.NoError(t, err)

	_, err = exec.Command("git", "push", "origin", gittest.DefaultBranch).CombinedOutput()
	require.NoError(t, err)

	log := gittest.LogRemote(t)
	require.Contains(t, log, "this commit is on the remote")
}

func TestLogRemoteDoesNotContainLocalCommits(t *testing.T) {
	gittest.InitRepository(t)

	_, err := exec.Command("git", "commit", "--allow-empty", "-m", "this commit is not on the remote").CombinedOutput()
	require.NoError(t, err)

	log := gittest.LogRemote(t)
	require.NotContains(t, log, "this commit is not on the remote")
}

func TestTagLocal(t *testing.T) {
	gittest.InitRepository(t)

	gittest.TagLocal(t, "0.1.0")

	out, err := exec.Command("git", "for-each-ref", "refs/tags").CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(out), "refs/tags/0.1.0")

	out, err = exec.Command("git", "ls-remote", "--tags").CombinedOutput()
	require.NoError(t, err)
	assert.NotContains(t, string(out), "refs/tags/0.1.0")
}

func TestShow(t *testing.T) {
	gittest.InitRepository(t)

	out := gittest.Show(t, gittest.DefaultBranch)
	assert.Contains(t, out, gittest.InitialCommit)
}

func TestCheckout(t *testing.T) {
	gittest.InitRepository(t)

	_, err := exec.Command("git", "branch", "testing").CombinedOutput()
	require.NoError(t, err)

	out := gittest.Checkout(t, "testing")
	assert.Equal(t, "Switched to branch 'testing'", out)
}

func TestRemote(t *testing.T) {
	gittest.InitRepository(t)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	remote := gittest.Remote(t)

	// Ensure path is sanitized before comparison
	assert.Equal(t, filepath.ToSlash(fmt.Sprintf("file://%s.git", cwd)), remote)
}
