package git_test

import (
	"fmt"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLocalCommits("testing git push"))

	client, _ := git.NewClient()
	out, err := client.Push()

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", gittest.DefaultBranch))

	remoteLog := gittest.RemoteLog(t)
	require.Equal(t, "testing git push", remoteLog[0].Message)
}

func TestPushWithPushOptions(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLocalCommits("testing git push options"))

	client, _ := git.NewClient()
	_, err := client.Push(git.WithPushOptions("option1", "option2"))

	require.NoError(t, err)
}

func TestPushResolveBranchError(t *testing.T) {
	nonWorkingDirectory(t)

	client, _ := git.NewClient()
	_, err := client.Push()

	assert.Error(t, err)
}

func TestPushAwareOfCurrentBranch(t *testing.T) {
	log := "(HEAD -> branch-aware, main, origin/main) chore: finished scaffolding project"
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithLocalCommits("this should be pushed on current branch"))

	client, _ := git.NewClient()
	out, err := client.Push()

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "branch-aware"))
}

func TestPushWithAllBranches(t *testing.T) {
	log := "(main, local-branch-1, local-branch-2) feat: can push all branches"
	gittest.InitRepository(t, gittest.WithLog(log))
	gittest.Tag(t, "0.1.0")

	client, _ := git.NewClient()
	out, err := client.Push(git.WithAllBranches())

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "main"))
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "local-branch-1"))
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "local-branch-2"))
	assert.NotContains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.1.0"))
}

func TestPushWithAllTags(t *testing.T) {
	log := "(main) feat: can push all tags"
	gittest.InitRepository(t, gittest.WithLog(log))
	gittest.Tag(t, "0.1.0")
	gittest.Tag(t, "0.2.0")

	client, _ := git.NewClient()
	out, err := client.Push(git.WithAllTags())

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.1.0"))
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.2.0"))
	assert.NotContains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "main"))
}

func TestPushWithRefSpecs(t *testing.T) {
	log := "(main, local-branch-3, local-branch-4) feat: can cherry-pick push"
	gittest.InitRepository(t, gittest.WithLog(log))
	gittest.Tag(t, "0.3.0")
	gittest.Tag(t, "0.4.0")

	client, _ := git.NewClient()
	out, err := client.Push(git.WithRefSpecs("0.3.0", "local-branch-3"))

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.3.0"))
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "local-branch-3"))
	assert.NotContains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.4.0"))
	assert.NotContains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "local-branch-4"))
}

func TestPushWithDeleteRefSpecs(t *testing.T) {
	log := "(tag: 0.1.0, tag: 0.2.0) feat: recreate user data indexes for speedier queries"
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	_, err := client.Push(git.WithDeleteRefSpecs("0.2.0"))
	require.NoError(t, err)

	remoteTags := gittest.RemoteTags(t)
	assert.ElementsMatch(t, []string{"0.1.0"}, remoteTags)
}
