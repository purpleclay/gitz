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
	"fmt"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLocalCommits("testing git push"))

	client, _ := git.NewClient()
	out, err := client.Push()

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", gittest.DefaultBranch))

	remoteLog := gittest.LogRemote(t)
	require.Equal(t, "testing git push", remoteLog[0].Commit)
}

func TestPushAwareOfCurrentBranch(t *testing.T) {
	log := "(HEAD -> branch-aware, main, origin/main) chore: finished scaffolding project"
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithLocalCommits("this should be pushed on current branch"))

	client, _ := git.NewClient()
	out, err := client.Push()

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "branch-aware"))
}

func TestPushWithAllBranches(t *testing.T) {
	log := "(main, local-branch-1, local-branch-2) feat: can push all branches"
	gittest.InitRepository(t, gittest.WithLog(log))
	gittest.TagLocal(t, "0.1.0")

	client, _ := git.NewClient()
	out, err := client.Push(git.WithAllBranches())

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "main"))
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "local-branch-1"))
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "local-branch-2"))
	assert.NotContains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.1.0"))
}

func TestPushWithAllTags(t *testing.T) {
	log := "(main) feat: can push all tags"
	gittest.InitRepository(t, gittest.WithLog(log))
	gittest.TagLocal(t, "0.1.0")
	gittest.TagLocal(t, "0.2.0")

	client, _ := git.NewClient()
	out, err := client.Push(git.WithAllTags())

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.1.0"))
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.2.0"))
	assert.NotContains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "main"))
}

func TestPushWithRefSpecs(t *testing.T) {
	log := "(main, local-branch-3, local-branch-4) feat: can cherry-pick push"
	gittest.InitRepository(t, gittest.WithLog(log))
	gittest.TagLocal(t, "0.3.0")
	gittest.TagLocal(t, "0.4.0")

	client, _ := git.NewClient()
	out, err := client.Push(git.WithRefSpecs("0.3.0", "local-branch-3"))

	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.3.0"))
	assert.Contains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "local-branch-3"))
	assert.NotContains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "0.4.0"))
	assert.NotContains(t, out, fmt.Sprintf("%[1]s -> %[1]s", "local-branch-4"))
}

func TestPushTag(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TagLocal(t, "0.1.0")

	client, _ := git.NewClient()
	_, err := client.PushTag("0.1.0")

	require.NoError(t, err)
	remoteTags := gittest.RemoteTags(t)
	assert.ElementsMatch(t, []string{"0.1.0"}, remoteTags)
}
