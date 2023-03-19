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
	"runtime"
	"strings"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Utility method for executing git commands and ensuring the trailing slash is trimmed
func gitExec(t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", args...).CombinedOutput()
	require.NoError(t, err)

	return strings.TrimSuffix(string(out), "\n")
}

func TestInitRepositoryConfigSet(t *testing.T) {
	gittest.InitRepository(t)

	cfg := gitExec(t, "config", "--list")

	assert.Contains(t, cfg, fmt.Sprintf("user.name=%s", gittest.DefaultAuthorName))
	assert.Contains(t, cfg, fmt.Sprintf("user.email=%s", gittest.DefaultAuthorEmail))
}

func TestInitRepositoryDefaultBranchSet(t *testing.T) {
	gittest.InitRepository(t)

	branch := gitExec(t, "branch", "--format=%(refname:short)")
	assert.Equal(t, gittest.DefaultBranch, branch)
}

func TestInitRepositoryWithLog(t *testing.T) {
	log := `chore: resolve broken build badge
ci: adopt new code security workflow`
	gittest.InitRepository(t, gittest.WithLog(log))

	out := gitExec(t, "log", "-n2", "--oneline")
	lines := strings.Split(out, "\n")
	require.Len(t, lines, 2)

	assert.Contains(t, lines[0], "chore: resolve broken build badge")
	assert.Contains(t, lines[1], "ci: adopt new code security workflow")
}

func TestInitRepositoryWithLogCreatesTags(t *testing.T) {
	log := `(tag: 0.1.0, tag: v1) feat: this is a brand new feature
ci: include github workflow`
	gittest.InitRepository(t, gittest.WithLog(log))

	localTags := localTags(t)
	remoteTags := remoteTags(t)

	assert.ElementsMatch(t, []string{"0.1.0", "v1"}, localTags)
	assert.ElementsMatch(t, []string{"0.1.0", "v1"}, remoteTags)
}

func localTags(t *testing.T) []string {
	t.Helper()
	tags := gitExec(t, "tag", "--format=%(refname:short)")

	return strings.Split(tags, "\n")
}

func remoteTags(t *testing.T) []string {
	t.Helper()
	tags := gitExec(t, "ls-remote", "--tags")

	cleanedTags := make([]string, 0)
	for _, tag := range strings.Split(tags, "\n") {
		if _, cleanedTag, found := strings.Cut(tag, "refs/tags/"); found {
			cleanedTags = append(cleanedTags, cleanedTag)
		}
	}

	return cleanedTags
}

func TestInitRepositoryWithLogCreatesBranches(t *testing.T) {
	log := `(main) chore: add example code snippets
(local-tracked) feat: support branch creation within log
(tracked, origin/tracked) docs: document fix
(origin/remote-tracked) fix: parsing of multiple tags within log
docs: update existing project README`
	gittest.InitRepository(t, gittest.WithLog(log))

	localBranches := localBranches(t)
	assert.Contains(t, localBranches, "local-tracked")
	assert.Contains(t, localBranches, "tracked")
	assert.NotContains(t, localBranches, "remote-tracked")

	remoteBranches := remoteBranches(t)
	assert.Contains(t, remoteBranches, "tracked")
	assert.Contains(t, remoteBranches, "remote-tracked")
	assert.NotContains(t, remoteBranches, "local-tracked")

	// Checkout and verify that branches are associated with the expected commit
	script := "git checkout $0 &>/dev/null; git log -n1 --oneline"
	out := shellExecInline(t, script, "local-tracked")
	assert.Contains(t, out, "feat: support branch creation within log")

	out = shellExecInline(t, script, "tracked")
	assert.Contains(t, out, "docs: document fix")

	out = shellExecInline(t, script, "remote-tracked")
	assert.Contains(t, out, "fix: parsing of multiple tags within log")
}

func shellExecInline(t *testing.T, inline string, args ...string) string {
	t.Helper()

	// Combine and squash args into a slice
	cmdArgs := append([]string{"-c", inline}, args...)

	interp := "/bin/bash"
	if runtime.GOOS == "windows" {
		interp = "bash"
	}

	out, err := exec.Command(interp, cmdArgs...).CombinedOutput()
	require.NoError(t, err)

	return string(out)
}

func localBranches(t *testing.T) []string {
	t.Helper()

	branches := gitExec(t, "branch", "--list", "--format=%(refname:short)")
	return strings.Split(branches, "\n")
}

func remoteBranches(t *testing.T) []string {
	t.Helper()

	branches := gitExec(t, "branch", "--list", "--remotes", "--format=%(refname:short)")

	cleanedBranches := make([]string, 0)
	for _, branch := range strings.Split(branches, "\n") {
		cleanedBranches = append(cleanedBranches, strings.TrimPrefix(branch, "origin/"))
	}
	return cleanedBranches
}

func TestInitRepositoryWithLogCheckoutBranch(t *testing.T) {
	log := `(HEAD -> branch-checkout, origin/branch-checkout) pass tests
write tests for branch checkout
(main, origin/main) docs: update existing project README`
	gittest.InitRepository(t, gittest.WithLog(log))

	currentBranch := gitExec(t, "branch", "--show-current")
	assert.Equal(t, "branch-checkout", currentBranch)

	out := gitExec(t, "log", "-n2", "--oneline")
	assert.Contains(t, out, "pass tests")
	assert.Contains(t, out, "write tests for branch checkout")
}

func TestInitRepositoryWithLogCheckoutBranchNotPushed(t *testing.T) {
	log := "(HEAD -> local-branch, main, origin/main) feat: latest and greatest feature"
	gittest.InitRepository(t, gittest.WithLog(log))

	currentBranch := gitExec(t, "branch", "--show-current")
	assert.Equal(t, "local-branch", currentBranch)

	remoteBranches := remoteBranches(t)
	assert.NotContains(t, remoteBranches, "local-branch")
}

func TestInitRepositoryWithFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("a.txt", "b.txt"))

	out := gitExec(t, "status", "--porcelain")
	status := strings.Split(out, "\n")

	assert.ElementsMatch(t, []string{"?? a.txt", "?? b.txt"}, status)
}

func TestInitRepositoryWithStagedFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles("c.txt", "dir/d.txt"))

	out := gitExec(t, "status", "--porcelain")
	status := strings.Split(out, "\n")

	assert.ElementsMatch(t, []string{"A  c.txt", "A  dir/d.txt"}, status)
}

func TestInitRepositoryWithLocalCommits(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLocalCommits("local commit 1", "local commit 2"))

	log := gitExec(t, "log", "--oneline")
	assert.Contains(t, log, "local commit 1")
	assert.Contains(t, log, "local commit 2")

	remoteLog := gitExec(t, "log", "--oneline", gittest.DefaultRemoteBranch)
	assert.NotContains(t, remoteLog, "local commit 1")
	assert.NotContains(t, remoteLog, "local commit 2")
}

func TestWithRemoteLog(t *testing.T) {
	log := "(main, origin/main) this is a remote commit"
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	localLog, err := exec.Command("git", "log", "-n1", "--oneline").CombinedOutput()
	require.NoError(t, err)
	assert.NotContains(t, string(localLog), "this is a remote commit")

	_, err = exec.Command("git", "pull").CombinedOutput()
	require.NoError(t, err)

	localLog, err = exec.Command("git", "log", "-n1", "--oneline").CombinedOutput()
	require.NoError(t, err)
	assert.Contains(t, string(localLog), "this is a remote commit")
}

func TestWithRemoteLogNewBranch(t *testing.T) {
	log := `(HEAD -> new-branch, origin/new-branch) pass tests
write tests for new feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	branches := remoteBranches(t)
	assert.NotContains(t, branches, "new-branch")

	gitExec(t, "pull")

	branches = remoteBranches(t)
	assert.Contains(t, branches, "new-branch")
}

func TestWithCloneDepth(t *testing.T) {
	log := `(main, origin/main) feat: this is commit number 3
feat: this is commit number 2
feat: this is commit number 1`

	gittest.InitRepository(t, gittest.WithLog(log), gittest.WithCloneDepth(1))

	out := gitExec(t, "log", "--oneline")
	lines := strings.Split(out, "\n")

	require.Len(t, lines, 1)
	assert.Contains(t, lines[0], "feat: this is commit number 3")
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

	gitExec(t, "tag", "0.1.0")
	gitExec(t, "tag", "0.2.0")

	tags := gittest.Tags(t)
	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0"}, tags)
}

func TestTagsIsEmpty(t *testing.T) {
	gittest.InitRepository(t)
	assert.Empty(t, gittest.Tags(t))
}

func TestRemoteTags(t *testing.T) {
	gittest.InitRepository(t)

	gitExec(t, "tag", "0.2.0")
	gitExec(t, "push", gittest.DefaultOrigin, "0.2.0")

	tags := gittest.RemoteTags(t)
	assert.ElementsMatch(t, []string{"0.2.0"}, tags)
}

func TestRemoteTagsIsEmpty(t *testing.T) {
	gittest.InitRepository(t)
	assert.Empty(t, gittest.RemoteTags(t))
}

func TestStageFile(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("test.txt"))

	gittest.StageFile(t, "test.txt")

	status := gitExec(t, "status", "--porcelain")
	assert.Contains(t, status, "A  test.txt")
}

func TestCommit(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles("file.txt"))

	gittest.Commit(t, "include file.txt")

	log := gitExec(t, "log", "-n1", "--oneline")
	assert.Contains(t, log, "include file.txt")
}

func TestLastCommit(t *testing.T) {
	gittest.InitRepository(t)

	gitExec(t, "commit", "--allow-empty", "-m", "this is a test")
	expectedHash := gitExec(t, "rev-parse", "HEAD")

	commit := gittest.LastCommit(t)
	assert.Equal(t, expectedHash, commit.Hash)
	assert.Equal(t, expectedHash[:7], commit.AbbrevHash)
	assert.Equal(t, gittest.DefaultAuthorName, commit.AuthorName)
	assert.Equal(t, gittest.DefaultAuthorEmail, commit.AuthorEmail)
	assert.Equal(t, "this is a test", commit.Message)
}

func TestPorcelainStatus(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("file1.txt", "file2.txt"))

	status := gittest.PorcelainStatus(t)
	assert.ElementsMatch(t, []string{"?? file1.txt", "?? file2.txt"}, status)
}

func TestProcelainStatusNoChanges(t *testing.T) {
	gittest.InitRepository(t)
	assert.Empty(t, gittest.PorcelainStatus(t))
}

func TestLogRemote(t *testing.T) {
	gittest.InitRepository(t)
	gitExec(t, "commit", "--allow-empty", "-m", "this commit is on the remote")
	gitExec(t, "push", "origin", gittest.DefaultBranch)

	log := gittest.LogRemote(t)

	require.Len(t, log, 2)
	require.Equal(t, "this commit is on the remote", log[0].Commit)
}

func TestRemoteLog(t *testing.T) {
	gittest.InitRepository(t)
	gitExec(t, "commit", "--allow-empty", "-m", "this commit is on the remote")
	gitExec(t, "push", "origin", gittest.DefaultBranch)

	log := gittest.RemoteLog(t)

	require.Len(t, log, 2)
	require.Equal(t, "this commit is on the remote", log[0].Commit)
}

func TestRemoteLogDoesNotContainLocalCommits(t *testing.T) {
	gittest.InitRepository(t)
	gitExec(t, "commit", "--allow-empty", "-m", "this commit is not on the remote")

	log := gittest.RemoteLog(t)

	require.Len(t, log, 1)
	assert.NotEqual(t, "this commit is not on the remote", log[0].Commit)
}

func TestLog(t *testing.T) {
	log := `(main, origin/main) chore: second line of the log
chore: first line of the log`
	gittest.InitRepository(t, gittest.WithLog(log))

	localLog := gittest.Log(t)
	require.Len(t, localLog, 3)
	assert.Equal(t, "chore: second line of the log", localLog[0].Message)
	assert.Equal(t, "chore: first line of the log", localLog[1].Message)
	assert.Equal(t, gittest.InitialCommit, localLog[2].Message)
}

func TestTagLocal(t *testing.T) {
	gittest.InitRepository(t)

	gittest.TagLocal(t, "0.1.0")

	localTags := localTags(t)
	assert.ElementsMatch(t, []string{"0.1.0"}, localTags)

	remoteTags := remoteTags(t)
	assert.Empty(t, remoteTags)
}

func TestShow(t *testing.T) {
	gittest.InitRepository(t)

	out := gittest.Show(t, gittest.DefaultBranch)
	assert.Contains(t, out, gittest.InitialCommit)
}

func TestCheckout(t *testing.T) {
	gittest.InitRepository(t)
	gitExec(t, "branch", "testing")

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

func TestShowBranch(t *testing.T) {
	gittest.InitRepository(t)

	branch := gittest.ShowBranch(t)
	assert.Equal(t, gittest.DefaultBranch, branch)
}

func TestBranches(t *testing.T) {
	gittest.InitRepository(t)

	script := `
for b in branch{1..3}; do
	git checkout -b $b;
done;`

	shellExecInline(t, script)

	branches := gittest.Branches(t)
	assert.ElementsMatch(t, []string{"branch1", "branch2", "branch3", gittest.DefaultBranch}, branches)
}

func TestBranchesOnInitializedRepository(t *testing.T) {
	changeToTmpDir(t)

	_, err := exec.Command("git", "init").CombinedOutput()
	require.NoError(t, err)

	branches := gittest.Branches(t)
	assert.Empty(t, branches)
}

func changeToTmpDir(t *testing.T) {
	t.Helper()
	changedFrom, err := os.Getwd()
	require.NoError(t, err)

	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(changedFrom))
	})
}

func TestRemoteBranches(t *testing.T) {
	gittest.InitRepository(t)

	script := `
for b in branch{1..3}; do
	git checkout -b $b;
done;
git push origin --all`

	shellExecInline(t, script)

	branches := gittest.RemoteBranches(t)
	assert.Contains(t, branches, "branch1")
	assert.Contains(t, branches, "branch2")
	assert.Contains(t, branches, "branch3")
}

func TestRemoteBranchesOnInitializedRepository(t *testing.T) {
	changeToTmpDir(t)

	_, err := exec.Command("git", "init").CombinedOutput()
	require.NoError(t, err)

	branches := gittest.RemoteBranches(t)
	assert.Empty(t, branches)
}
