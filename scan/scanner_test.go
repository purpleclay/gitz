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

package scan_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/purpleclay/gitz/scan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrefixedLines(t *testing.T) {
	text := `> this is line #1
>this is line #2
>    this is line #3
and it is spread over two lines   `

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(scan.PrefixedLines('>'))

	lines := readUntilEOF(t, scanner)
	require.Len(t, lines, 3)
	assert.Equal(t, "this is line #1", lines[0])
	assert.Equal(t, "this is line #2", lines[1])
	assert.Equal(t, `this is line #3
and it is spread over two lines`, lines[2])
}

func TestPrefixedLinesIgnoresNonLeadingPrefix(t *testing.T) {
	text := "this was created by jdoe >>> <jdoe@test.com>"

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(scan.PrefixedLines('>'))

	lines := readUntilEOF(t, scanner)
	require.Len(t, lines, 1)
	assert.Equal(t, "this was created by jdoe >>> <jdoe@test.com>", lines[0])
}

func readUntilEOF(t *testing.T, scanner *bufio.Scanner) []string {
	t.Helper()

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func TestPrefixedLinesNoPrefix(t *testing.T) {
	text := `this is line #1
this is line #2
this is line #3
and it is spread over two lines
`

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(scan.PrefixedLines('>'))

	lines := readUntilEOF(t, scanner)
	require.Len(t, lines, 1)
	assert.Equal(t, `this is line #1
this is line #2
this is line #3
and it is spread over two lines`, lines[0])
}

func TestPrefixedLinesInconsistentPrefixUse(t *testing.T) {
	text := `this is line #1
this is line #2
> this is line #3
and it is spread over two lines`

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(scan.PrefixedLines('>'))

	lines := readUntilEOF(t, scanner)
	require.Len(t, lines, 2)
	assert.Equal(t, `this is line #1
this is line #2`, lines[0])
	assert.Equal(t, `this is line #3
and it is spread over two lines`, lines[1])
}
