package git_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
)

func TestClone(t *testing.T) {
	log := "(main, origin/main) chore: testing if a git clone works"
	gittest.InitRepository(t, gittest.WithLog(log))

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote)
	require.NoError(t, err)

	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))
	assert.Equal(t, "chore: testing if a git clone works", gittest.LastCommit(t).Message)
}

func TestCloneWithDirectory(t *testing.T) {
	gittest.InitRepository(t)

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote, git.WithDirectory("cloned-repo"))

	require.NoError(t, err)
	assert.NoDirExists(t, gittest.ClonedRepositoryName)
	assert.DirExists(t, "cloned-repo")
}

func TestCloneWithDepth(t *testing.T) {
	log := `(main, origin/main) chore: testing clone depth line 2
chore: testing clone depth line 1`
	gittest.InitRepository(t, gittest.WithLog(log))

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote, git.WithDepth(1))

	require.NoError(t, err)
	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))

	localLog := gittest.Log(t)
	require.Len(t, localLog, 1)
	assert.Equal(t, "chore: testing clone depth line 2", localLog[0].Message)
}

func TestCloneWithDepthLessThanOne(t *testing.T) {
	log := `(main, origin/main) chore: testing clone depth line 2
chore: testing clone depth line 1`
	gittest.InitRepository(t, gittest.WithLog(log))

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote, git.WithDepth(0))

	require.NoError(t, err)
	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))

	localLog := gittest.Log(t)
	assert.Len(t, localLog, 3)
}

func TestCloneWithCheckoutRefForBranch(t *testing.T) {
	log := "(main, origin/main, origin/branch-cloning) chore: test branch is cloned"
	gittest.InitRepository(t, gittest.WithLog(log))

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote, git.WithCheckoutRef("branch-cloning"))

	require.NoError(t, err)
	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))
	assert.Equal(t, "branch-cloning", gittest.ShowBranch(t))
}

func TestCloneWithCheckoutRefForTag(t *testing.T) {
	log := `(main, origin/main) chore: shouldn't see this commit
(tag: clone-tag) chore: test this tag is cloned`
	gittest.InitRepository(t, gittest.WithLog(log))

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote, git.WithCheckoutRef("clone-tag"))

	require.NoError(t, err)
	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))
	assert.Equal(t, "chore: test this tag is cloned", gittest.LastCommit(t).Message)
}

func TestCloneWithCheckoutRefEmptyString(t *testing.T) {
	log := `(main, origin/main) chore: shouldn't see this commit
(tag: clone-tag) chore: test this tag is cloned`
	gittest.InitRepository(t, gittest.WithLog(log))

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote, git.WithCheckoutRef("   "))

	require.NoError(t, err)
	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))
	assert.Equal(t, "chore: shouldn't see this commit", gittest.LastCommit(t).Message)
}

func TestCloneWithNoTags(t *testing.T) {
	log := "(main, origin/main) chore: test no tags are cloned"
	gittest.InitRepository(t, gittest.WithLog(log))

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote, git.WithNoTags())

	require.NoError(t, err)
	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))
	assert.Empty(t, gittest.Tags(t))
}
