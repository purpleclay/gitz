package git_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
)

func TestPull(t *testing.T) {
	log := "(tag: 0.1.0, main, origin/main) feat: a new exciting feature"
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	require.NotEqual(t, "feat: a new exciting feature", gittest.LastCommit(t).Message)

	client, _ := git.NewClient()
	_, err := client.Pull()
	require.NoError(t, err)

	assert.Equal(t, "feat: a new exciting feature", gittest.LastCommit(t).Message)
	tags := gittest.Tags(t)
	assert.ElementsMatch(t, []string{"0.1.0"}, tags)
}

func TestPullWithFetchIgnoreTags(t *testing.T) {
	log := `(tag: 0.3.0, main, origin/main) feat: third feature
(tag: 0.2.0) feat: second feature
(tag: 0.1.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	client, _ := git.NewClient()
	_, err := client.Pull(git.WithFetchIgnoreTags())
	require.NoError(t, err)

	assert.Empty(t, gittest.Tags(t))
}

func TestPullWithPullRefSpecs(t *testing.T) {
	log := `(main, origin/main) test: add test for validating refspecs
(origin/branch) fix: ensure pull supports refspecs`
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	_, err := client.Pull(git.WithPullRefSpecs("branch:branch1"))

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"main", "branch1"}, gittest.Branches(t))
}
