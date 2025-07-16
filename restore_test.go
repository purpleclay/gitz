package git_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
)

func TestRestoreUsingForUntrackedFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("main.go", ".github/ci.yaml", "go.mod"))

	untracked := [2]git.FileStatusIndicator{git.Untracked, git.Untracked}

	client, _ := git.NewClient()
	err := client.RestoreUsing([]git.FileStatus{
		{Indicators: untracked, Path: "main.go"},
		{Indicators: untracked, Path: ".github/"},
		{Indicators: untracked, Path: "go.mod"},
	})
	require.NoError(t, err)

	statuses := gittest.PorcelainStatus(t)
	assert.Empty(t, statuses)
}

func TestRestoreUsingForModifiedFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithCommittedFiles("main.go", "doc.go"))
	gittest.WriteFile(t, "main.go", "updated", 0o644)
	gittest.WriteFile(t, "doc.go", "updated", 0o644)
	gittest.StageFile(t, "main.go")

	client, _ := git.NewClient()
	err := client.RestoreUsing([]git.FileStatus{
		{Indicators: [2]git.FileStatusIndicator{git.Modified, git.Untracked}, Path: "main.go"},
		{Indicators: [2]git.FileStatusIndicator{git.Untracked, git.Modified}, Path: "doc.go"},
	})
	require.NoError(t, err)

	statuses := gittest.PorcelainStatus(t)
	assert.Empty(t, statuses)
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
