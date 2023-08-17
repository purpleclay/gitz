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

func TestPullWithFetchNoTags(t *testing.T) {
	log := `(tag: 0.3.0, main, origin/main) feat: third feature
(tag: 0.2.0) feat: second feature
(tag: 0.1.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))

	client, _ := git.NewClient()
	_, err := client.Pull(git.WithFetchNoTags())
	require.NoError(t, err)

	assert.Empty(t, gittest.Tags(t))
}

func TestPullWithFetchAllTags(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TagRemote(t, "0.1.0")
	gittest.TagRemote(t, "0.2.0")
	require.Empty(t, gittest.Tags(t))

	client, _ := git.NewClient()
	_, err := client.Pull(git.WithFetchAllTags())
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0"}, gittest.Tags(t))
}

func TestPullWithFetchDepth(t *testing.T) {
	log := `(main, origin/main) feat: third feature
feat: second feature
feat: first feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))
	shallowClone(t, gittest.Remote(t))

	client, _ := git.NewClient()
	_, err := client.Pull(git.WithFetchDepth(2))
	require.NoError(t, err)

	glog := gittest.Log(t)
	assert.Len(t, glog, 2)
}

func shallowClone(t *testing.T, remote string) {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))
	gittest.MustExec(t, "git clone --depth=1 -- "+remote)

	require.NoError(t, os.Chdir(gittest.ClonedRepositoryName))
}

func TestPullWithFetchRefSpecs(t *testing.T) {
	log := `(main, origin/main) test: add test for validating refspecs
(origin/branch) fix: ensure pull supports refspecs`
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	_, err := client.Pull(git.WithFetchRefSpecs("branch:branch1"))

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"main", "branch1"}, gittest.Branches(t))
}
