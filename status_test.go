package git_test

import (
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPorcelainStatus(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("README.md"), gittest.WithStagedFiles("go.mod"))

	client, _ := git.NewClient()
	statuses, err := client.PorcelainStatus()
	require.NoError(t, err)

	require.Len(t, statuses, 2)
	assert.ElementsMatch(t,
		[]string{"?? README.md", "A  go.mod"},
		[]string{statuses[0].String(), statuses[1].String()})
}

func TestClean(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	clean, err := client.Clean()
	require.NoError(t, err)

	assert.True(t, clean)
}

func TestCleanWithStagedChanges(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles("example.txt"))

	client, _ := git.NewClient()
	clean, err := client.Clean()
	require.NoError(t, err)

	assert.False(t, clean)
}
