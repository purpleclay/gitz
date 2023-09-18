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
	"strings"
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	_, err := client.Tag("0.1.0")

	require.NoError(t, err)

	localTags := gittest.Tags(t)
	assert.ElementsMatch(t, []string{"0.1.0"}, localTags)

	remoteTags := gittest.RemoteTags(t)
	assert.ElementsMatch(t, []string{"0.1.0"}, remoteTags)
}

func TestTagWithInvalidName(t *testing.T) {
	gittest.InitRepository(t)

	// See https://git-scm.com/docs/git-check-ref-format for details on what
	// constitutes an invalid tag (ref)
	client, _ := git.NewClient()
	_, err := client.Tag("[0.1.0]")

	assert.ErrorContains(t, err, "'[0.1.0]' is not a valid tag name")
}

func TestTagWithAnnotation(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	_, err := client.Tag("0.1.0", git.WithAnnotation("created tag 0.1.0"))

	require.NoError(t, err)

	out := gittest.Show(t, "0.1.0")
	assert.Contains(t, out, fmt.Sprintf("Tagger: %s", gittest.DefaultAuthorLog))
	assert.Contains(t, out, "created tag 0.1.0")
}

func TestTagWithSkipSigning(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.signingkey", "DOES-NOT-EXIST", "tag.gpgsign", "true")

	client, _ := git.NewClient()
	_, err := client.Tag("0.1.0", git.WithSkipSigning())

	require.NoError(t, err)
}

func TestTagWithTagConfig(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	_, err := client.Tag("0.1.0",
		git.WithAnnotation("test inline git options"),
		git.WithTagConfig("user.name", "bane", "user.email", "bane@dc.com"))

	require.NoError(t, err)
	out := gittest.Show(t, "0.1.0")
	assert.Contains(t, out, "Tagger: bane <bane@dc.com>")
}

func TestTagWithLocalOnly(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	_, err := client.Tag("0.1.0", git.WithLocalOnly())

	require.NoError(t, err)
	assert.Empty(t, gittest.RemoteTags(t))
}

func TestTagWithCommitRef(t *testing.T) {
	log := `ci: add extra job to workflow for running golden file tests
test: expand current test suite using golden files`
	gittest.InitRepository(t, gittest.WithLog(log))
	glog := gittest.Log(t)
	require.Len(t, glog, 3)

	client, _ := git.NewClient()
	_, err := client.Tag("0.1.1", git.WithCommitRef(glog[1].Hash))

	require.NoError(t, err)
	out := gittest.Show(t, "0.1.1")
	assert.Contains(t, out, "commit "+glog[1].Hash)
}

func TestDeleteTags(t *testing.T) {
	log := "(tag: 0.1.0, tag: 0.2.0) feat(ui): add new fancy button to ui"
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	_, err := client.DeleteTags([]string{"0.1.0", "0.2.0"})
	require.NoError(t, err)

	localTags := gittest.Tags(t)
	assert.Empty(t, localTags)

	remoteTags := gittest.RemoteTags(t)
	assert.Empty(t, remoteTags)
}

func TestDeleteTagsLocally(t *testing.T) {
	log := "(tag: 0.1.0, tag: 0.2.0) fix: indexed data is not in the correct order"
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	_, err := client.DeleteTags([]string{"0.1.0", "0.2.0"}, git.WithLocalDelete())
	require.NoError(t, err)

	localTags := gittest.Tags(t)
	assert.Empty(t, localTags)

	remoteTags := gittest.RemoteTags(t)
	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0"}, remoteTags)
}

func TestTags(t *testing.T) {
	log := `(tag: 0.2.0, tag: v1) feat: add support for tag sorting and filtering
(tag: 0.1.0) feat: add support for basic cloning`
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	tags, err := client.Tags()

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0", "v1"}, tags)
}

func TestTagsEmpty(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	tags, err := client.Tags()

	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestTagsWithShellGlob(t *testing.T) {
	log := `(tag: 0.2.0, tag: v1) feat: add support for tag sorting and filtering
(tag: 0.1.0) feat: add support for basic cloning`
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	tags, err := client.Tags(git.WithShellGlob("*.*.*"))

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0"}, tags)
}

func TestTagsWithSortBy(t *testing.T) {
	log := `(tag: 0.11.0) feat: add support for tag sorting and filtering
(tag: 0.10.0) feat: add support for inspecting a repository
(tag: 0.9.1) fix: grep pattern not working as expected
(tag: 0.9.0) feat: add support for log filtering`
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	tags, err := client.Tags(git.WithSortBy(git.CreatorDateDesc, git.VersionDesc))

	require.NoError(t, err)
	require.Len(t, tags, 4)
	assert.Equal(t, "0.11.0", tags[0])
	assert.Equal(t, "0.10.0", tags[1])
	assert.Equal(t, "0.9.1", tags[2])
	assert.Equal(t, "0.9.0", tags[3])
}

func TestTagsWithSortBySemanticVersions(t *testing.T) {
	log := "(tag: 0.1.0, tag: 0.2.0-beta.1, tag: 0.2.0-beta.2, tag: 0.2.0) fix: support semantic version sorting"
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	tags, err := client.Tags(git.WithSortBy(git.VersionDesc))

	require.NoError(t, err)
	require.Len(t, tags, 4)
	assert.Equal(t, "0.2.0", tags[0])
	assert.Equal(t, "0.2.0-beta.2", tags[1])
	assert.Equal(t, "0.2.0-beta.1", tags[2])
	assert.Equal(t, "0.1.0", tags[3])
}

func TestTagsQueryingTagsError(t *testing.T) {
	nonWorkingDirectory(t)

	client, _ := git.NewClient()
	_, err := client.Tags()

	require.Error(t, err)
}

func TestTagsWithCount(t *testing.T) {
	log := "(tag: 0.1.0, tag: 0.2.0, tag: 0.3.0, tag: 0.4.0) feat: limit tag retrieval"
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	tags, err := client.Tags(git.WithCount(3))

	require.NoError(t, err)
	require.Len(t, tags, 3)
	assert.Equal(t, "0.1.0", tags[0])
	assert.Equal(t, "0.2.0", tags[1])
	assert.Equal(t, "0.3.0", tags[2])
}

func TestTagsWithCountEqualToMax(t *testing.T) {
	log := "(tag: 0.1.0, tag: 0.2.0) feat: limit tag retrieval"
	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	tags, err := client.Tags(git.WithCount(2))

	require.NoError(t, err)
	require.Len(t, tags, 2)
	assert.Equal(t, "0.1.0", tags[0])
	assert.Equal(t, "0.2.0", tags[1])
}

func TestTagsWithFilters(t *testing.T) {
	log := `(tag: ui/0.2.0, tag: ui/v1) feat: replace text within table with pills
(tag: backend/0.2.0, tag: backend/v1) feat: support sorting of items through api
(tag: ui/0.1.0) feat: add paging support on results table
(tag: backend/0.1.0) feat: extend api to return item states`
	gittest.InitRepository(t, gittest.WithLog(log))

	uiFilter := func(tag string) bool {
		return strings.HasPrefix(tag, "ui/")
	}

	noVTagsFilter := func(tag string) bool {
		return !strings.HasSuffix(tag, "v1")
	}

	client, _ := git.NewClient()
	tags, err := client.Tags(git.WithFilters(uiFilter, noVTagsFilter, nil))

	require.NoError(t, err)
	require.Len(t, tags, 2)
	assert.Equal(t, "ui/0.1.0", tags[0])
	assert.Equal(t, "ui/0.2.0", tags[1])
}
