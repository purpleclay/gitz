package git

import (
	"fmt"
	"strings"
)

// RestoreUsing will restore a given set of files back to their previous
// known state within the current repository (working directory). By
// inspecting each files [FileStatus], the correct decision can be made
// when restoring it.
func (c *Client) RestoreUsing(statuses []FileStatus) error {
	for _, status := range statuses {
		var err error

		switch {
		case status.Untracked():
			err = c.removeUntrackedFile(status.Path)
		case status.Modified():
			err = c.restoreFile(status)
		case status.Renamed():
			err = c.undoRenamedFile(status.Path)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) removeUntrackedFile(pathspec string) error {
	_, err := c.Exec("git clean --force -- " + pathspec)
	return err
}

func (c *Client) restoreFile(status FileStatus) error {
	var buf strings.Builder
	buf.WriteString("git restore ")
	if status.Indicators[0] == Modified {
		buf.WriteString("--staged --worktree ")
	}
	buf.WriteString(status.Path)

	_, err := c.Exec(buf.String())
	return err
}

func (c *Client) undoRenamedFile(pathspec string) error {
	original, renamed, _ := strings.Cut(pathspec, porcelainRenameSeparator)

	_, err := c.Exec(fmt.Sprintf("git mv %s %s", renamed, original))
	return err
}
