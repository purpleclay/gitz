package git_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
)

func TestShowBlobs(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, "README.md", "The quick brown fox jumps over the lazy dog")
	gittest.StageFile(t, "README.md")
	gittest.TempFile(t, "LICENSE", "Follow the yellow brick road")
	gittest.StageFile(t, "LICENSE")
	gittest.Commit(t, "docs: include readme")

	// Lookup refs
	readmeRef := gittest.ObjectRef(t, "README.md")
	licenseRef := gittest.ObjectRef(t, "LICENSE")

	client, _ := git.NewClient()
	blobs, err := client.ShowBlobs(readmeRef, licenseRef)
	require.NoError(t, err)

	require.Len(t, blobs, 2)
	assert.Equal(t, readmeRef, blobs[readmeRef].Ref)
	assert.Equal(t, "The quick brown fox jumps over the lazy dog", blobs[readmeRef].Diff)

	assert.Equal(t, licenseRef, blobs[licenseRef].Ref)
	assert.Equal(t, "Follow the yellow brick road", blobs[licenseRef].Diff)
}

func TestShowTrees(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles(
		"README.md",
		"pkg/scan/scanner.go",
		"pkg/scan/scanner_test.go",
		"internal/task/scan.go",
		"internal/task.go",
		"internal/tui/dashboard.go"))
	gittest.Commit(t, "chore: commit everything")

	// Lookup refs
	pkgRef := gittest.ObjectRef(t, "pkg")
	internalRef := gittest.ObjectRef(t, "internal")

	client, _ := git.NewClient()
	trees, err := client.ShowTrees(pkgRef, internalRef)
	require.NoError(t, err)

	require.Len(t, trees, 2)
	assert.Equal(t, pkgRef, trees[pkgRef].Ref)
	assert.ElementsMatch(t, []string{"scan/"}, trees[pkgRef].Entries)

	assert.Equal(t, internalRef, trees[internalRef].Ref)
	assert.ElementsMatch(t, []string{"task.go", "task/", "tui/"}, trees[internalRef].Entries)
}

func TestShowCommits(t *testing.T) {
	gittest.InitRepository(t)
	gittest.CommitEmptyWithAuthor(t, "joker", "joker@dc.com", "docs: document new parsing features")
	gittest.CommitEmpty(t, `feat: add functionality to parse a commit

ensure a commit can be parsed when using the git show command`)

	entries := gittest.Log(t)

	client, _ := git.NewClient()
	commits, err := client.ShowCommits(entries[0].Hash, entries[1].Hash)
	require.NoError(t, err)

	require.Len(t, commits, 2)
	ref := entries[0].Hash
	assert.Equal(t, ref, commits[ref].Ref)
	assert.Equal(t, "batman", commits[ref].Author.Name)
	assert.Equal(t, "batman@dc.com", commits[ref].Author.Email)
	assert.WithinDuration(t, time.Now(), commits[ref].AuthorDate, time.Second*2)
	assert.Equal(t, "batman", commits[ref].Committer.Name)
	assert.Equal(t, "batman@dc.com", commits[ref].Committer.Email)
	assert.WithinDuration(t, time.Now(), commits[ref].CommitterDate, time.Second*2)
	assert.Equal(t, `feat: add functionality to parse a commit

ensure a commit can be parsed when using the git show command`, commits[ref].Message)

	ref = entries[1].Hash
	assert.Equal(t, ref, commits[ref].Ref)
	assert.Equal(t, "joker", commits[ref].Author.Name)
	assert.Equal(t, "joker@dc.com", commits[ref].Author.Email)
	assert.WithinDuration(t, time.Now(), commits[ref].AuthorDate, time.Second*2)
	assert.Equal(t, "batman", commits[ref].Committer.Name)
	assert.Equal(t, "batman@dc.com", commits[ref].Committer.Email)
	assert.WithinDuration(t, time.Now(), commits[ref].CommitterDate, time.Second*2)
	assert.Equal(t, "docs: document new parsing features", commits[ref].Message)
}

func TestShowTags(t *testing.T) {
	gittest.InitRepository(t)
	gittest.Tag(t, "0.1.0")
	gittest.TagAnnotated(t, "0.2.0", "chore: tagged release at 0.2.0")

	client, _ := git.NewClient()
	tags, err := client.ShowTags("0.1.0", "0.2.0")
	require.NoError(t, err)

	require.Len(t, tags, 2)
	tag := tags["0.1.0"]
	assert.Equal(t, "0.1.0", tag.Ref)
	assert.Nil(t, tag.Annotation)
	assert.Equal(t, gittest.InitialCommit, tag.Commit.Message)

	tag = tags["0.2.0"]
	assert.Equal(t, "0.2.0", tag.Ref)
	require.NotNil(t, tag.Annotation)
	assert.Equal(t, "batman", tag.Annotation.Tagger.Name)
	assert.Equal(t, "batman@dc.com", tag.Annotation.Tagger.Email)
	assert.WithinDuration(t, time.Now(), tag.Annotation.TaggerDate, time.Second*2)
	assert.Equal(t, "chore: tagged release at 0.2.0", tag.Annotation.Message)
	assert.Equal(t, gittest.InitialCommit, tag.Commit.Message)
}
