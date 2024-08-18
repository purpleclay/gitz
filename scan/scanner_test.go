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

func TestDiffLines(t *testing.T) {
	text := `diff --git a/clone.go b/clone.go
index f181e5f..bea7426 100644
--- a/clone.go
+++ b/clone.go
@@ -10,6 +10,7 @@ import (
 // repository is cloned onto the file system into a target working directory
 type CloneOption func(*cloneOptions)

+// Hello
 type cloneOptions struct {
        Config      []string
        CheckoutRef string
diff --git a/commit.go b/commit.go
index 906a132..2e6954c 100644
--- a/commit.go
+++ b/commit.go
@@ -10,6 +10,7 @@ import (
 // created against the current repository (working directory)
 type CommitOption func(*commitOptions)

+// Hello, again!
 type commitOptions struct {
        AllowEmpty    bool
        Config        []string
`

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(scan.DiffLines())

	lines := readUntilEOF(t, scanner)
	require.Len(t, lines, 2)
	assert.Equal(t, `diff --git a/clone.go b/clone.go
index f181e5f..bea7426 100644
--- a/clone.go
+++ b/clone.go
@@ -10,6 +10,7 @@ import (
 // repository is cloned onto the file system into a target working directory
 type CloneOption func(*cloneOptions)

+// Hello
 type cloneOptions struct {
        Config      []string
        CheckoutRef string`, lines[0])
	assert.Equal(t, `diff --git a/commit.go b/commit.go
index 906a132..2e6954c 100644
--- a/commit.go
+++ b/commit.go
@@ -10,6 +10,7 @@ import (
 // created against the current repository (working directory)
 type CommitOption func(*commitOptions)

+// Hello, again!
 type commitOptions struct {
        AllowEmpty    bool
        Config        []string`, lines[1])
}
