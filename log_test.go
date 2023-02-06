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

func TestLog(t *testing.T) {
	t.Parallel()
	log := `fix: parsing error when input string is too long
ci: extend the existing build workflow to include integration tests
docs: create initial mkdocs material documentation
feat: add second operation to library
feat: add first operation to library`

	gittest.InitRepository(t, gittest.WithLog(log))

	client, _ := git.NewClient()
	out, err := client.Log()

	require.NoError(t, err)
	assert.Contains(t, out, "fix: parsing error when input string is too long")
	assert.Contains(t, out, "ci: extend the existing build workflow to include integration tests")
	assert.Contains(t, out, "docs: create initial mkdocs material documentation")
	assert.Contains(t, out, "feat: add second operation to library")
	assert.Contains(t, out, "feat: add first operation to library")
}
