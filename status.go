package git

import (
	"bufio"
	"fmt"
	"strings"
)

// FileStatusIndicator contains a single character that represents
// a files status within a git repository. Based on the git
// specification: https://git-scm.com/docs/git-status#_output.
type FileStatusIndicator byte

const (
	// Added indicates a new file that has been added to the index (staged).
	Added FileStatusIndicator = 'A'

	// Copied indicates a file that has been copied from another file,
	// with the copy being tracked by git.
	Copied FileStatusIndicator = 'C'

	// Deleted indicates a file that has been deleted from the working tree
	// and the deletion has been staged.
	Deleted FileStatusIndicator = 'D'

	// Ignored indicates a file that is being ignored by git due to
	// .gitignore rules or other ignore mechanisms.
	Ignored FileStatusIndicator = '!'

	// Modified indicates a file that has been modified in the working tree
	// and/or index compared to HEAD.
	Modified FileStatusIndicator = 'M'

	// Renamed indicates a file that has been renamed, with git detecting
	// the rename operation.
	Renamed FileStatusIndicator = 'R'

	// TypeChanged indicates a file whose type has changed (e.g., regular file
	// to symlink, or vice versa).
	TypeChanged FileStatusIndicator = 'T'

	// Updated indicates a file with merge conflicts that have been resolved
	// and staged, but not yet committed (unmerged -> staged).
	Updated FileStatusIndicator = 'U'

	// Unmodified indicates a file that has not been changed (clean state).
	Unmodified FileStatusIndicator = ' '

	// Untracked indicates a file that is not being tracked by git.
	Untracked FileStatusIndicator = '?'
)

const porcelainRenameSeparator = " -> "

// FileStatus represents the status of a file within a repository.
type FileStatus struct {
	// Indicators is a two character array that contains
	// the current status of a file within both the current index
	// and the working repository tree.
	//
	// Examples:
	//
	// 	'??' - a file that is not tracked
	//	' A' - a file that has been added to the working tree
	//	'M ' - a file that has been modified within the index
	Indicators [2]FileStatusIndicator

	// Path of the file relative to the root of the
	// current repository
	Path string
}

// String representation of a file status that adheres to the
// porcelain v1 format.
func (f FileStatus) String() string {
	return fmt.Sprintf("%c%c %s", f.Indicators[0], f.Indicators[1], f.Path)
}

// Untracked identifies whether a file is not currently tracked.
func (f FileStatus) Untracked() bool {
	return f.Indicators[0] == Untracked && f.Indicators[1] == Untracked
}

// Modified idenfities whether a file has been modified and therefore contains changes.
func (f FileStatus) Modified() bool {
	return f.Indicators[0] == Modified || f.Indicators[1] == Modified
}

// Renamed identifies whether a file has been renamed.
func (f FileStatus) Renamed() bool {
	return f.Indicators[0] == Renamed
}

// StatusOption provides a way for setting specific options during a
// porcelain status operation. Each support option can customize the list
// of file statuses identified within the current repository (working directory).
type StatusOption func(*statusOptions)

type statusOptions struct {
	IgnoreRenames   bool
	IgnoreUntracked bool
}

// WithIgnoreRenames will turn off rename detection, removing any renamed
// files or directories from the retrieved file statuses.
func WithIgnoreRenames() StatusOption {
	return func(opts *statusOptions) {
		opts.IgnoreRenames = true
	}
}

// WithIgnoreUntracked will remove any untracked files from the retrieved
// file statuses.
func WithIgnoreUntracked() StatusOption {
	return func(opts *statusOptions) {
		opts.IgnoreUntracked = true
	}
}

// PorcelainStatus identifies if there are any changes within the current
// repository (working directory) and returns them in the parseable
// porcelain v1 format.
func (c *Client) PorcelainStatus(opts ...StatusOption) ([]FileStatus, error) {
	options := &statusOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var buf strings.Builder
	buf.WriteString("git status --porcelain")

	if options.IgnoreRenames {
		buf.WriteString(" --no-renames")
	}

	if options.IgnoreUntracked {
		buf.WriteString(" --untracked-files=no")
	}

	log, err := c.Exec(buf.String())
	if err != nil {
		return nil, err
	}

	return parsePorcelainV1(log), nil
}

// Clean determines if the current repository (working directory) is in
// a clean state. A repository is deemed clean, if it contains no changes.
func (c *Client) Clean() (bool, error) {
	statuses, err := c.PorcelainStatus()
	return len(statuses) == 0, err
}

func parsePorcelainV1(log string) []FileStatus {
	var statuses []FileStatus

	scanner := bufio.NewScanner(strings.NewReader(log))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		statuses = append(statuses, FileStatus{
			Indicators: [2]FileStatusIndicator{
				FileStatusIndicator(line[0]),
				FileStatusIndicator(line[1]),
			},
			Path: line[3:],
		})
	}

	return statuses
}
