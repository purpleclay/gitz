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
