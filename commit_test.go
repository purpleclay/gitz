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
