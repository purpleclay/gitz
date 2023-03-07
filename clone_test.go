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
	"os"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClone(t *testing.T) {
	log := "chore: testing if a git clone works"
	gittest.InitRepository(t, gittest.WithLog(log))

	// Grab the remote for cloning later
	remote := gittest.Remote(t)

	// Clone the existing repository into a new temporary directory
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	client, _ := git.NewClient()
	_, err := client.Clone(remote)
	require.NoError(t, err)

	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))
	lastCommit := gittest.LastCommit(t)
	assert.Contains(t, lastCommit, log)
}
