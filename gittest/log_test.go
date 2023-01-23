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
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLog(t *testing.T) {
	log := `(tag: 0.1.0) feat: improve existing cli documentation
docs: create initial mkdocs material documentation
feat: add secondary cli command to support filtering of results
feat: scaffold initial cli and add first command`

	entries := gittest.ParseLog(log)

	require.Len(t, entries, 4)
	assert.Equal(t, entries[0].Commit, "feat: improve existing cli documentation")
	assert.Equal(t, entries[0].Tag, "0.1.0")
	assert.Equal(t, entries[1].Commit, "docs: create initial mkdocs material documentation")
	assert.Empty(t, entries[1].Tag)
	assert.Equal(t, entries[2].Commit, "feat: add secondary cli command to support filtering of results")
	assert.Empty(t, entries[2].Tag)
	assert.Equal(t, entries[3].Commit, "feat: scaffold initial cli and add first command")
	assert.Empty(t, entries[3].Tag)
}

func TestParseLogTrimsSpaces(t *testing.T) {
	log := "   feat: testing if leading and trailing spaces are removed   "

	entries := gittest.ParseLog(log)

	require.Len(t, entries, 1)
	assert.Equal(t, entries[0].Commit, "feat: testing if leading and trailing spaces are removed")
}

func TestParseLogMalformedCommit(t *testing.T) {
	log := "(not a tag) feat: this raw commit should be returned"

	entries := gittest.ParseLog(log)

	require.Len(t, entries, 1)
	assert.Equal(t, entries[0].Commit, "(not a tag) feat: this raw commit should be returned")
}
