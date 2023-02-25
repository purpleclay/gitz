/*
Copyright (c) 2023 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package git_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStage(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("file.txt", "dir1/file.txt", "dir2/file.txt"))

	client, _ := git.NewClient()
	_, err := client.Stage()

	require.NoError(t, err)
	status := gittest.PorcelainStatus(t)

	statusLines := parsePorcelainStatus(t, status)
	require.Len(t, statusLines, 3)

	assert.ElementsMatch(t, statusLines, []string{
		"A  file.txt",
		"A  dir1/file.txt",
		"A  dir2/file.txt",
	})
}

func parsePorcelainStatus(t *testing.T, status string) []string {
	t.Helper()

	scanner := bufio.NewScanner(strings.NewReader(status))
	scanner.Split(bufio.ScanLines)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
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
	fmt.Println(status)

	statusLines := parsePorcelainStatus(t, status)
	require.Len(t, statusLines, 4)

	assert.ElementsMatch(t, statusLines, []string{
		"A  file.txt",
		"?? dir1/file.txt",
		"A  dir1/file.gif",
		"?? dir2/",
	})
}

func TestStageWithPathSpecsIgnoresEmptyPathSpecs(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("file1.txt", "file2.txt"))

	client, _ := git.NewClient()
	_, err := client.Stage(git.WithPathSpecs(" ", "   file2.txt   "))

	require.NoError(t, err)
	status := gittest.PorcelainStatus(t)
	fmt.Println(status)

	statusLines := parsePorcelainStatus(t, status)
	require.Len(t, statusLines, 2)

	assert.ElementsMatch(t, statusLines, []string{"?? file1.txt", "A  file2.txt"})
}
