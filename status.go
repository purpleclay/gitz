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

// PorcelainStatus identifies if there are any changes within the current
// repository (working directory) and returns them in the parseable
// porcelain v1 format
func (c *Client) PorcelainStatus() ([]FileStatus, error) {
	log, err := c.exec("git status --porcelain")
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
