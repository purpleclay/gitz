package git

import (
	"fmt"
	"strings"
)

// RestoreUsing ...(against HEAD only)
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
		buf.WriteString("--staged ")
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
