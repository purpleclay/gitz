package git_test

import (
	"os"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	log := `(origin/branch1) feat: extend the use of filtering when searching
(origin/branch2) fix: broken ordering of filters`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	client, _ := git.NewClient()
	_, err := client.Fetch()
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"main"}, gittest.Branches(t))
	assert.ElementsMatch(t, []string{
		"branch1",
		"branch2",
		gittest.DefaultBranch,
		gittest.DefaultOrigin,
	}, gittest.RemoteBranches(t))
}

func TestFetchWithIgnoreTags(t *testing.T) {
	log := `(tag: 0.3.0, main, origin/main) feat: third feature
(tag: 0.2.0) feat: second feature
(tag: 0.1.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	client, _ := git.NewClient()
	_, err := client.Fetch(git.WithIgnoreTags())
	require.NoError(t, err)

	assert.Empty(t, gittest.Tags(t))
	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0", "0.3.0"}, gittest.RemoteTags(t))
}

func TestFetchWithTags(t *testing.T) {
	log := `(tag: 0.3.0, main, origin/main) feat: third feature
(tag: 0.2.0) feat: second feature
(tag: 0.1.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	client, _ := git.NewClient()
	_, err := client.Fetch(git.WithTags())
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0", "0.3.0"}, gittest.Tags(t))
}

func TestFetchWithDepthTo(t *testing.T) {
	log := `(main, origin/main) feat: third feature
feat: second feature
feat: first feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))
	shallowClone(t, gittest.Remote(t))

	client, _ := git.NewClient()
	_, err := client.Fetch(git.WithDepthTo(2))
	require.NoError(t, err)

	glog := gittest.Log(t)
	assert.Len(t, glog, 2)
}

func shallowClone(t *testing.T, remote string) {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))
	gittest.MustExec(t, "git clone --depth=1 -- "+remote)

	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))
}

func TestFetchWithFetchRefSpecs(t *testing.T) {
	log := "(main) feat: ensure fetch supports refspecs"
	remLog := "(origin/main) test: add test for validating refspecs"
	gittest.InitRepository(t, gittest.WithLog(log), gittest.WithRemoteLog(remLog))

	client, _ := git.NewClient()
	_, err := client.Fetch(git.WithFetchRefSpecs("main"))
	require.NoError(t, err)

	dlog := gittest.LogBetween(t, "main", "origin/main")
	require.Len(t, dlog, 1)
	assert.Equal(t, "test: add test for validating refspecs", dlog[0].Message)
}

func TestFetchWithUnshallow(t *testing.T) {
	log := `(main, origin/main) fifth feature
feat: fourth feature
feat: third feature
feat: second feature
feat: first feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))
	shallowClone(t, gittest.Remote(t))

	client, _ := git.NewClient()
	_, err := client.Fetch(git.WithUnshallow())
	require.NoError(t, err)

	glog := gittest.Log(t)
	assert.Len(t, glog, 6)
}
