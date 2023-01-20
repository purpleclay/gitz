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
	gittest.InitRepo(t)

	client := git.NewClient()
	err := client.Tag("0.1.0")

	require.NoError(t, err)

	out := gittest.Tags(t)
	assert.Contains(t, out, refs("0.1.0"))

	out = gittest.RemoteTags(t)
	assert.Contains(t, out, refs("0.1.0"))
}

func TestDeleteTag(t *testing.T) {
	log := "(tag: 0.1.0) feat: a brand new feature"

	gittest.InitRepo(t, gittest.WithLog(log))

	client := git.NewClient()
	err := client.DeleteTag("0.11.0")

	require.NoError(t, err)

	out := gittest.Tags(t)
	assert.NotContains(t, out, refs("0.1.0"))

	out = gittest.RemoteTags(t)
	assert.NotContains(t, out, refs("0.1.0"))
}
