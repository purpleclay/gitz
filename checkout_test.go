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
