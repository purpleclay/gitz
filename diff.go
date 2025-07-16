package git

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/purpleclay/chomp"

	"github.com/purpleclay/gitz/scan"
)

const (
	// git diff header delimiter > @@ ... @@.
	hdrDelim = "@@"
	// prefix for lines added.
	addPrefix = "+"
	// prefix for lines removed.
	remPrefix = "-"
)

// DiffOption provides a way for setting specific options during a diff
// operation. Each supported option can customize the way the diff is
// executed against the current repository (working directory).
type DiffOption func(*diffOptions)

type diffOptions struct {
	DiffPaths []string
}

// WithDiffPaths allows the diff to be targeted to specific files and
// folders within the current repository (working directory). Paths to
// files and folders are relative to the root of the repository. All
// leading and trailing whitespace will be trimmed from the file paths,
// allowing empty paths to be ignored.
func WithDiffPaths(paths ...string) DiffOption {
	return func(opts *diffOptions) {
		opts.DiffPaths = trim(paths...)
	}
}

// FileDiff represents a snapshot containing all of the changes to
// a file within a repository (working directory).
type FileDiff struct {
	// Path of the file within the repository (working directory)
	Path string

	// DiffChunk contains all of the identified changes within
	// the file
	Chunks []DiffChunk
}

// DiffChunk represents a snapshot of a single change (chunk) to
// a file within a repository (working directory).
type DiffChunk struct {
	// Added optionally contains details of the text that has
	// been added to a file as part of the current change
	Added DiffChange

	// Removed optionally contains details of the text that has
	// been removed from a file as part of the current change
	Removed DiffChange
}

// DiffChange captures details about an individual chunk
// within a git diff. It contains both the changed text and
// its exact position (and line count) within the file.
type DiffChange struct {
	// LineNo is the position within the file where the
	// change starts
	LineNo int

	// Count is the number of lines that has changed
	Count int

	// Change contains the text that has changed
	Change string
}

// Diff captures the changes made to files within the current repository (working
// directory). Options can be provided to customize how the current diff is
// determined. By default, all diffs (or changes) to files within the repository
// will be retrieved. The diff is generated using the following git options:
//
//	git diff -U0 --no-color
func (c *Client) Diff(opts ...DiffOption) ([]FileDiff, error) {
	options := &diffOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var buf strings.Builder
	buf.WriteString("git diff -U0 --no-color")

	if len(options.DiffPaths) > 0 {
		buf.WriteString(" -- ")
		buf.WriteString(strings.Join(options.DiffPaths, " "))
	}

	out, err := c.Exec(buf.String())
	if err != nil {
		return nil, err
	}
	return parseDiffs(out)
}

func parseDiffs(log string) ([]FileDiff, error) {
	var diffs []FileDiff

	scanner := bufio.NewScanner(strings.NewReader(log))
	scanner.Split(scan.DiffLines())

	for scanner.Scan() {
		diff, err := parseDiff(scanner.Text())
		if err != nil {
			return nil, err
		}

		diffs = append(diffs, diff)
	}

	return diffs, nil
}

func parseDiff(diff string) (FileDiff, error) {
	rem, path, err := diffPath()(diff)
	if err != nil {
		return FileDiff{}, err
	}

	rem, _, err = chomp.Until(hdrDelim)(rem)
	if err != nil {
		return FileDiff{}, err
	}

	chunks, err := diffChunks(rem)
	if err != nil {
		return FileDiff{}, err
	}

	return FileDiff{
		Path:   path,
		Chunks: chunks,
	}, nil
}

func diffPath() chomp.Combinator[string] {
	return func(s string) (string, string, error) {
		var rem string
		var err error

		if rem, _, err = chomp.Tag("diff --git ")(s); err != nil {
			return rem, "", err
		}

		var path string
		if rem, path, err = chomp.Until(" ")(rem); err != nil {
			return rem, "", err
		}
		path = path[strings.Index(path, "/")+1:]

		rem, _, err = chomp.Eol()(rem)
		return rem, path, err
	}
}

func diffChunks(in string) ([]DiffChunk, error) {
	_, chunks, err := chomp.Map(chomp.Many(diffChunk()),
		func(in []string) []DiffChunk {
			var diffChunks []DiffChunk

			for i := 0; i+5 < len(in); i += 6 {
				chunk := DiffChunk{
					Removed: DiffChange{
						LineNo: mustInt(in[i]),
						Count:  mustInt(in[i+1]),
						Change: in[i+4],
					},
					Added: DiffChange{
						LineNo: mustInt(in[i+2]),
						Count:  mustInt(in[i+3]),
						Change: in[i+5],
					},
				}

				if chunk.Added.Count == 0 {
					chunk.Added.Count = 1
				}

				if chunk.Removed.Count == 0 {
					chunk.Removed.Count = 1
				}

				diffChunks = append(diffChunks, chunk)
			}

			return diffChunks
		},
	)(in)

	return chunks, err
}

func mustInt(in string) int {
	out, _ := strconv.Atoi(in)
	return out
}

func diffChunk() chomp.Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var rem string
		var err error

		var changes []string
		rem, changes, err = chomp.Delimited(
			chomp.Tag(hdrDelim+" "),
			chomp.SepPair(diffChunkHeaderChange(remPrefix), chomp.Tag(" "), diffChunkHeaderChange(addPrefix)),
			chomp.Eol(),
		)(s)
		if err != nil {
			return rem, nil, err
		}

		var removed string
		rem, removed, err = chomp.Map(
			chomp.ManyN(chomp.Prefixed(chomp.Eol(), chomp.Tag(remPrefix)), 0),
			func(in []string) string { return strings.Join(in, "\n") },
		)(rem)
		if err != nil {
			return rem, nil, err
		}

		var added string
		rem, added, err = chomp.Map(
			chomp.ManyN(chomp.Prefixed(chomp.Eol(), chomp.Tag(addPrefix)), 0),
			func(in []string) string { return strings.Join(in, "\n") },
		)(rem)
		if err != nil {
			return rem, nil, err
		}

		return rem, append(changes, removed, added), nil
	}
}

func diffChunkHeaderChange(prefix string) chomp.Combinator[[]string] {
	return func(s string) (string, []string, error) {
		rem, _, err := chomp.Tag(prefix)(s)
		if err != nil {
			return rem, nil, err
		}

		return chomp.All(
			chomp.While(chomp.IsDigit),
			chomp.Opt(chomp.Prefixed(chomp.While(chomp.IsDigit), chomp.Tag(","))),
		)(rem)
	}
}
