package git_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
)

func TestStage(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("file.txt", "dir1/file.txt", "dir2/file.txt"))

	client, _ := git.NewClient()
	_, err := client.Stage()

	require.NoError(t, err)
	status := gittest.PorcelainStatus(t)
	assert.ElementsMatch(t, []string{
		"A  file.txt",
		"A  dir1/file.txt",
		"A  dir2/file.txt",
	}, status)
}

func TestStageWithPathSpecs(t *testing.T) {
	files := []string{
		"file.txt",
		"dir1/file.txt",
		"dir1/file.gif",
		"dir2/file.txt",
	}
	gittest.InitRepository(t, gittest.WithFiles(files...))

	client, _ := git.NewClient()
	_, err := client.Stage(git.WithPathSpecs("file.txt", "dir1/*.gif"))

	require.NoError(t, err)
	status := gittest.PorcelainStatus(t)
	assert.ElementsMatch(t, []string{
		"A  file.txt",
		"?? dir1/file.txt",
		"A  dir1/file.gif",
		"?? dir2/",
	}, status)
}

func TestStageWithPathSpecsIgnoresEmptyPathSpecs(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("file1.txt", "file2.txt"))

	client, _ := git.NewClient()
	_, err := client.Stage(git.WithPathSpecs(" ", "   file2.txt   "))

	require.NoError(t, err)
	status := gittest.PorcelainStatus(t)
	assert.ElementsMatch(t, []string{"?? file1.txt", "A  file2.txt"}, status)
}

func TestStaged(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("README.md"),
		gittest.WithStagedFiles("go.mod", "pkg/config/config.go"))

	client, _ := git.NewClient()
	staged, err := client.Staged()
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"go.mod", "pkg/config/config.go"}, staged)
}
