package git

import (
	"bufio"
	"fmt"
	"strings"
)

// FileStatusIndicator contains a single character that represents
// a files status within a git repository. Based on the git
// specification: https://git-scm.com/docs/git-status#_output
type FileStatusIndicator byte

const (
	Added       FileStatusIndicator = 'A'
	Copied      FileStatusIndicator = 'C'
	Deleted     FileStatusIndicator = 'D'
	Ignored     FileStatusIndicator = '!'
	Modified    FileStatusIndicator = 'M'
	Renamed     FileStatusIndicator = 'R'
	TypeChanged FileStatusIndicator = 'T'
	Updated     FileStatusIndicator = 'U'
	Unmodified  FileStatusIndicator = ' '
	Untracked   FileStatusIndicator = '?'
)

const porcelainRenameSeparator = " -> "

// FileStatus represents the status of a file within a repository
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
// porcelain v1 format
func (f FileStatus) String() string {
	return fmt.Sprintf("%c%c %s", f.Indicators[0], f.Indicators[1], f.Path)
}

// Untracked identifies whether a file is not currently tracked
func (f FileStatus) Untracked() bool {
	return f.Indicators[0] == Untracked && f.Indicators[1] == Untracked
}

// Modified idenfities whether a file has been modified and therefore contains changes
func (f FileStatus) Modified() bool {
	return f.Indicators[0] == Modified || f.Indicators[1] == Modified
}

// Renamed identifies whether a file has been renamed
func (f FileStatus) Renamed() bool {
	return f.Indicators[0] == Renamed
}

// StatusOption provides a way for setting specific options during a
// porcelain status operation. Each support option can customize the list
// of file statuses identified within the current repository (working directory)
type StatusOption func(*statusOptions)

type statusOptions struct {
	IgnoreRenames   bool
	IgnoreUntracked bool
}

// WithIgnoreRenames will turn off rename detection, removing any renamed
// files or directories from the retrieved file statuses
func WithIgnoreRenames() StatusOption {
	return func(opts *statusOptions) {
		opts.IgnoreRenames = true
	}
}

// WithIgnoreUntracked will remove any untracked files from the retrieved
// file statuses
func WithIgnoreUntracked() StatusOption {
	return func(opts *statusOptions) {
		opts.IgnoreUntracked = true
	}
}

// PorcelainStatus identifies if there are any changes within the current
// repository (working directory) and returns them in the parseable
// porcelain v1 format
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

	log, err := c.exec(buf.String())
	if err != nil {
		return nil, err
	}

	return parsePorcelainV1(log), nil
}

// Clean determines if the current repository (working directory) is in
// a clean state. A repository is deemed clean, if it contains no changes
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
