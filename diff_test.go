package git_test

import (
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	gittest.InitRepository(t,
		gittest.WithCommittedFiles("main.go"),
		gittest.WithFileContent("main.go", `package main

import "fmt"

func print() {
	fmt.Println("Hello, World!")
}

func main() {
	print()
}`))

	overwriteFile(t, "main.go", `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Hello, %s\n" + os.Args[1])
}`)

	client, _ := git.NewClient()
	diffs, err := client.Diff()
	require.NoError(t, err)

	require.Len(t, diffs, 1)
	assert.Equal(t, "main.go", diffs[0].Path)

	require.Len(t, diffs[0].Chunks, 2)
	assert.Equal(t, 3, diffs[0].Chunks[0].Added.LineNo)
	assert.Equal(t, 4, diffs[0].Chunks[0].Added.Count)
	assert.Equal(t, `import (
	"fmt"
	"os"
)`, diffs[0].Chunks[0].Added.Change)

	assert.Equal(t, 3, diffs[0].Chunks[0].Removed.LineNo)
	assert.Equal(t, 5, diffs[0].Chunks[0].Removed.Count)
	assert.Equal(t, `import "fmt"

func print() {
	fmt.Println("Hello, World!")
}`, diffs[0].Chunks[0].Removed.Change)

	assert.Equal(t, 9, diffs[0].Chunks[1].Added.LineNo)
	assert.Equal(t, 1, diffs[0].Chunks[1].Added.Count)
	assert.Equal(t, `	fmt.Printf("Hello, %s\n" + os.Args[1])`, diffs[0].Chunks[1].Added.Change)

	assert.Equal(t, 10, diffs[0].Chunks[1].Removed.LineNo)
	assert.Equal(t, 1, diffs[0].Chunks[1].Removed.Count)
	assert.Equal(t, `	print()`, diffs[0].Chunks[1].Removed.Change)
}
