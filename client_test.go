package git_test

import (
	"fmt"
	"path/filepath"
	"strings"
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
	assert.Equal(t, gittest.WorkingDirectory(t), repo.RootDir)
	assert.Equal(t, repo.Origin, gittest.Remote(t))
	require.Len(t, repo.Remotes, 1)
	assert.Equal(t, repo.Remotes[gittest.DefaultOrigin], gittest.Remote(t))
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

	hash := gittest.LastCommit(t).Hash
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

func TestRepositoryWithMultipleRemotes(t *testing.T) {
	gittest.InitRepository(t)
	gittest.Exec(t, "git remote add gitlab git@gitlab.com:purpleclay/test.git")

	client, _ := git.NewClient()
	repo, err := client.Repository()
	require.NoError(t, err)

	require.Len(t, repo.Remotes, 2)
	assert.Equal(t, repo.Remotes[gittest.DefaultOrigin], gittest.Remote(t))
	assert.Equal(t, repo.Remotes["gitlab"], "git@gitlab.com:purpleclay/test.git")
}

func TestToRelativePath(t *testing.T) {
	gittest.InitRepository(t)
	root := gittest.WorkingDirectory(t)

	client, _ := git.NewClient()
	rel, err := client.ToRelativePath(filepath.Join(root, "a/nested/directory"))

	require.NoError(t, err)
	assert.Equal(t, "a/nested/directory", rel)
}

func TestToRelativePathNotInWorkingDirectoryError(t *testing.T) {
	gittest.InitRepository(t)
	root := gittest.WorkingDirectory(t)
	// ensure it is agnostic to the OS
	rel := osDriveLetter(t, root) + "/a/non/related/path"

	client, _ := git.NewClient()
	_, err := client.ToRelativePath(rel)

	// Cope with unwiedly paths due to temporary test directories
	assert.EqualError(t, err,
		fmt.Sprintf("%s is not relative to the git repository working directory %s as it produces path %s",
			rel, root, makeRelativeTo(t, rel, root)))
}

func osDriveLetter(t *testing.T, path string) string {
	t.Helper()
	return path[0:strings.Index(path, "/")]
}

func makeRelativeTo(t *testing.T, path, target string) string {
	t.Helper()
	n := strings.Count(target, "/")

	// Remove any drive letter
	relPath := strings.TrimPrefix(path, osDriveLetter(t, path))
	relPath = strings.TrimPrefix(relPath, "/")
	return strings.Repeat("../", n) + relPath
}
