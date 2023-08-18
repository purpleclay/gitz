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

func TestPull(t *testing.T) {
	log := "(tag: 0.1.0, main, origin/main) feat: a new exciting feature"
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	require.NotEqual(t, gittest.LastCommit(t).Message, "feat: a new exciting feature")

	client, _ := git.NewClient()
	_, err := client.Pull()
	require.NoError(t, err)

	assert.Equal(t, gittest.LastCommit(t).Message, "feat: a new exciting feature")
	tags := gittest.Tags(t)
	assert.ElementsMatch(t, []string{"0.1.0"}, tags)
}

func TestPullWithFetchIgnoreTags(t *testing.T) {
	log := `(tag: 0.3.0, main, origin/main) feat: third feature
(tag: 0.2.0) feat: second feature
(tag: 0.1.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	client, _ := git.NewClient()
	_, err := client.Pull(git.WithFetchIgnoreTags())
	require.NoError(t, err)

	assert.Empty(t, gittest.Tags(t))
}

func TestPullWithPullRefSpecs(t *testing.T) {
	log := `(main, origin/main) test: add test for validating refspecs
(origin/branch) fix: ensure pull supports refspecs`
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	_, err := client.Pull(git.WithPullRefSpecs("branch:branch1"))

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"main", "branch1"}, gittest.Branches(t))
}
