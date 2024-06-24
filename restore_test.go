package git_test

import (
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestoreUsingForUntrackedFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("README.md", ".github/ci.yaml", "go.mod"))

	untracked := [2]git.FileStatusIndicator{git.Untracked, git.Untracked}

	client, _ := git.NewClient()
	err := client.RestoreUsing([]git.FileStatus{
		{Indicators: untracked, Path: "README.md"},
		{Indicators: untracked, Path: ".github/"},
		{Indicators: untracked, Path: "go.mod"},
	})
	require.NoError(t, err)

	statuses := gittest.PorcelainStatus(t)
	assert.Empty(t, statuses)
}

func TestRestoreUsingForModifiedFiles(t *testing.T) {
	// TODO: committed files
	// TODO: modify one and stage
	// TODO: modify the other and do not stage
}

func TestRestoreUsingForRenamedFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithCommittedFiles("main.go", "cache.go", "keys.go"))
	gittest.Move(t, "cache.go", "internal/cache/cache.go")
	gittest.Move(t, "keys.go", "internal/cache/keys.go")

	renamed := [2]git.FileStatusIndicator{git.Renamed, git.Unmodified}

	client, _ := git.NewClient()
	err := client.RestoreUsing([]git.FileStatus{
		{Indicators: renamed, Path: "cache.go -> internal/cache/cache.go"},
		{Indicators: renamed, Path: "keys.go -> internal/cache/keys.go"},
	})
	require.NoError(t, err)

	statuses := gittest.PorcelainStatus(t)
	assert.Empty(t, statuses)
}
