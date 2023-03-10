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

// Formats a tag into the expected refs/tags/<tag> format
func refs(tag string) string {
	return fmt.Sprintf("refs/tags/%s", tag)
}

func TestTag(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	_, err := client.Tag("0.1.0")

	require.NoError(t, err)

	out := gittest.Tags(t)
	assert.Contains(t, out, refs("0.1.0"))

	out = gittest.RemoteTags(t)
	assert.Contains(t, out, refs("0.1.0"))
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

func TestTagWithAnnotationIgnores(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "EmptyString",
			message: "",
		},
		{
			name:    "StringWithOnlyWhitespace",
			message: "     ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gittest.InitRepository(t)

			client, _ := git.NewClient()
			_, err := client.Tag("0.1.0", git.WithAnnotation(tt.message))

			require.NoError(t, err)

			out := gittest.Show(t, "0.1.0")
			assert.NotContains(t, out, fmt.Sprintf("Tagger: %s", gittest.DefaultAuthorLog))
		})
	}
}

func TestDeleteTag(t *testing.T) {
	log := "(tag: 0.1.0) feat: a brand new feature"

	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	_, err := client.DeleteTag("0.1.0")

	require.NoError(t, err)

	out := gittest.Tags(t)
	assert.NotContains(t, out, refs("0.1.0"))

	out = gittest.RemoteTags(t)
	assert.NotContains(t, out, refs("0.1.0"))
}

func TestDeleteMissingLocalTag(t *testing.T) {
	gittest.InitRepository(t)

	client, _ := git.NewClient()
	_, err := client.DeleteTag("0.1.0")

	assert.ErrorContains(t, err, "tag '0.1.0' not found")
}
