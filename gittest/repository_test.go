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

package gittest_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRepositoryConfigSet(t *testing.T) {
	gittest.InitRepo(t)

	out, err := exec.Command("git", "config", "--list").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), fmt.Sprintf("user.name=%s", gittest.DefaultAuthorName))
	assert.Contains(t, string(out), fmt.Sprintf("user.email=%s", gittest.DefaultAuthorName))
}

func TestInitRepositoryDefaultBranchSet(t *testing.T) {
	gittest.InitRepo(t)

	out, err := exec.Command("git", "branch").CombinedOutput()
	require.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("* %s\n", gittest.DefaultBranch), string(out))
}

func TestInitRepositoryWithLog(t *testing.T) {
	log := "feat: this is a brand new feature"
	gittest.InitRepo(t, gittest.WithLog(log))

	out, err := exec.Command("git", "log", "--oneline").CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), "feat: this is a brand new feature")
}

func TestExecHasRawGitOutput(t *testing.T) {
	out := gittest.Exec(t, "git --version")

	assert.Contains(t, out, "git version")
}

func TestTags(t *testing.T) {
	gittest.InitRepo(t)

	_, err := exec.Command("git", "tag", "0.1.0").CombinedOutput()
	require.NoError(t, err)

	out := gittest.Tags(t)
	assert.Contains(t, out, "refs/tags/0.1.0")
}

func TestRemoteTags(t *testing.T) {
	gittest.InitRepo(t)

	_, err := exec.Command("git", "tag", "0.2.0").CombinedOutput()
	require.NoError(t, err)

	_, err = exec.Command("git", "push", "origin", "0.2.0").CombinedOutput()
	require.NoError(t, err)

	out := gittest.RemoteTags(t)
	assert.Contains(t, out, "refs/tags/0.2.0")
}
