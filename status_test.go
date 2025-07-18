package git_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
)

func TestPorcelainStatus(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("main.go"), gittest.WithStagedFiles("go.mod"))

	client, _ := git.NewClient()
	statuses, err := client.PorcelainStatus()
	require.NoError(t, err)

	require.Len(t, statuses, 2)
	assert.ElementsMatch(t,
		[]string{"?? main.go", "A  go.mod"},
		[]string{statuses[0].String(), statuses[1].String()},
	)
}

func TestPorcelainStatusWithIgnoreUntracked(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("main.go"), gittest.WithStagedFiles("go.mod"))

	client, _ := git.NewClient()
	statuses, err := client.PorcelainStatus(git.WithIgnoreUntracked())
	require.NoError(t, err)

	require.Len(t, statuses, 1)
	assert.ElementsMatch(t, []string{"A  go.mod"}, []string{statuses[0].String()})
}

func TestPorcelainStatusWithIgnoreRenames(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("go.mod"))
	gittest.Move(t, "README.md", "CONTRIBUTING.md")

	client, _ := git.NewClient()
	statuses, err := client.PorcelainStatus(git.WithIgnoreRenames())
	require.NoError(t, err)

	require.Len(t, statuses, 3)
	assert.ElementsMatch(t,
		[]string{"?? go.mod", "A  CONTRIBUTING.md", "D  README.md"},
		[]string{statuses[0].String(), statuses[1].String(), statuses[2].String()},
	)
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
