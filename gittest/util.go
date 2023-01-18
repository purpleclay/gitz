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

package gittest

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// InitRepo ...
func InitRepo(t *testing.T) {
	t.Helper()

	// Track our current directory
	current, err := os.Getwd()
	require.NoError(t, err)

	// Generate two temporary directories. The first is initialized as a
	// bare repository and becomes our filesystem based remote. The second
	// is our working repository, which is a clone of the former
	tmpDir := t.TempDir()
	require.NoError(t, os.Chdir(tmpDir))

	Exec(t, "git init --bare test.git")
	Exec(t, "git clone ./test.git")

	require.NoError(t, os.Chdir("./test"))

	// Initialize the repository so that is ready for use
	Exec(t, "git commit --allow-empty -m 'initialize repository'")

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(current))
	})
}

// Exec ...
func Exec(t *testing.T, cmd string) string {
	t.Helper()

	p, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	require.NoError(t, err)

	var buf bytes.Buffer
	r, err := interp.New(
		interp.StdIO(os.Stdin, &buf, &buf),
	)
	require.NoError(t, err)
	require.NoError(t, r.Run(context.Background(), p))

	return buf.String()
}
