package git

import (
	"fmt"
	"strings"
)

// RestoreUsing will restore a given set of files back to their previous
// known state within the current repository (working directory). By
// inspecting each files [FileStatus], the correct decision can be made
// when restoring it
func (c *Client) RestoreUsing(statuses []FileStatus) error {
	for _, status := range statuses {
		var err error

		if status.Untracked() {
			err = c.removeUntrackedFile(status.Path)
		} else if status.Modified() {
			err = c.restoreFile(status)
		} else if status.Renamed() {
			err = c.undoRenamedFile(status.Path)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) removeUntrackedFile(pathspec string) error {
	_, err := c.exec("git clean --force -- " + pathspec)
	return err
}

func (c *Client) restoreFile(status FileStatus) error {
	var buf strings.Builder
	buf.WriteString("git restore ")
	if status.Indicators[0] == Modified {
		buf.WriteString("--staged --worktree ")
	}
	buf.WriteString(status.Path)

	_, err := c.exec(buf.String())
	return err
}

func (c *Client) undoRenamedFile(pathspec string) error {
	original, renamed, _ := strings.Cut(pathspec, porcelainRenameSeparator)

	_, err := c.exec(fmt.Sprintf("git mv %s %s", renamed, original))
	return err
}
