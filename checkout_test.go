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

func TestCheckout(t *testing.T) {
	log := `(HEAD -> branch-checkout, origin/branch-checkout) pass tests
write tests for branch checkout
(main, origin/main) docs: update existing project README`
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	out, err := client.Checkout("main")
	require.NoError(t, err)

	// Inspect the raw git output
	assert.Contains(t, out, "Switched to branch 'main'")
	assert.Equal(t, gittest.LastCommit(t).Message, "docs: update existing project README")
}

func TestCheckoutCreatesLocalBranch(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	_, err := client.Checkout("testing")
	require.NoError(t, err)

	branches := gittest.Branches(t)
	remoteBranches := gittest.RemoteBranches(t)

	assert.Contains(t, branches, "testing")
	assert.NotContains(t, remoteBranches, "testing")
}

func TestCheckoutQueryingBranchesError(t *testing.T) {
	nonWorkingDirectory(t)

	client, _ := git.NewClient()
	_, err := client.Checkout("testing")

	require.Error(t, err)
}
