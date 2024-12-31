package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

const (
	disabledNumericOption = -1

	// RelativeAtRoot can be used to compare if a path is equivalent to the
	// root of a current git repository working directory
	RelativeAtRoot = "."

	// HeadRef is a pointer to the latest commit within a git repository
	HeadRef = "HEAD"
)

// ErrGitMissing is raised when no git client was identified
// within the PATH environment variable on the current OS
type ErrGitMissing struct {
	// PathEnv contains the value of the PATH environment variable
	PathEnv string
}

// Error returns a friendly formatted message of the current error
func (e ErrGitMissing) Error() string {
	return fmt.Sprintf("git is not installed under the PATH environment variable. PATH resolves to %s", e.PathEnv)
}

// ErrGitExecCommand is raised when a git command fails to execute
type ErrGitExecCommand struct {
	// Cmd contains the command that caused the git client to error
	Cmd string

	// Out contains any raw output from the git client as a result
	// of the error
	Out string
}

// Error returns a friendly formatted message of the current error
func (e ErrGitExecCommand) Error() string {
	return fmt.Sprintf(`failed to execute git command: %s

%s`, e.Cmd, e.Out)
}

// ErrGitNonRelativePath is raised when attempting to resolve a path
// within a git repository that isn't relative to the root of the
// working directory
type ErrGitNonRelativePath struct {
	// RootDir contains the root working directory of the repository
	RootDir string

	// TargetPath contains the path that was resolved against the
	// root working directory of the repository
	TargetPath string

	// RelativePath contains the resolved relative path which raised
	// the error
	RelativePath string
}

// Error returns a friendly formatted message of the current error
func (e ErrGitNonRelativePath) Error() string {
	return fmt.Sprintf("%s is not relative to the git repository working directory %s as it produces path %s",
		e.TargetPath, e.RootDir, e.RelativePath)
}

// Repository provides a snapshot of the current state of a repository
// (working directory)
type Repository struct {
	// DetachedHead is true if the current repository HEAD points to a
	// specific commit, rather than a branch
	DetachedHead bool

	// DefaultBranch is the initial branch that is checked out when
	// a repository is cloned
	DefaultBranch string

	// Origin contains the URL of the remote which this repository
	// was cloned from
	Origin string

	// Remotes will contain all of the remotes and their URLs as
	// configured for this repository
	Remotes map[string]string

	// RootDir contains the path to the cloned directory
	RootDir string

	// ShallowClone is true if the current repository has been cloned
	// to a specified depth without the entire commit history
	ShallowClone bool
}

// Client provides a way of performing fluent operations against git.
// Any git operation exposed by this client are effectively handed-off
// to an installed git client on the current OS. Git operations will be
// mapped as closely as possible to the official Git specification
type Client struct {
	gitVersion string
}

// NewClient returns a new instance of the git client
func NewClient() (*Client, error) {
	c := &Client{}

	if _, err := c.exec("type git"); err != nil {
		return nil, ErrGitMissing{PathEnv: os.Getenv("PATH")}
	}

	c.gitVersion, _ = c.exec("git --version")
	return c, nil
}

// Version of git used by the client
func (c *Client) Version() string {
	return c.gitVersion
}

// Repository captures and returns a snapshot of the current repository
// (working directory) state
func (c *Client) Repository() (Repository, error) {
	isRepo, _ := c.exec("git rev-parse --is-inside-work-tree")
	if strings.TrimSpace(isRepo) != "true" {
		return Repository{}, errors.New("current working directory is not a git repository")
	}

	isShallow, _ := c.exec("git rev-parse --is-shallow-repository")
	isDetached, _ := c.exec("git branch --show-current")
	defaultBranch, _ := c.exec("git rev-parse --abbrev-ref remotes/origin/HEAD")
	rootDir, _ := c.rootDir()

	// Identify all remotes associated with this repository. If this is a new
	// locally initialized repository, this could be empty
	rmts, _ := c.exec("git remote")
	remotes := map[string]string{}
	for _, remote := range strings.Split(rmts, "\n") {
		remoteURL, _ := c.exec("git remote get-url " + remote)
		remotes[remote] = filepath.ToSlash(remoteURL)
	}

	origin := ""
	if orig, found := remotes["origin"]; found {
		origin = orig
	}

	return Repository{
		DetachedHead:  strings.TrimSpace(isDetached) == "",
		DefaultBranch: strings.TrimPrefix(defaultBranch, "origin/"),
		Origin:        origin,
		Remotes:       remotes,
		RootDir:       rootDir,
		ShallowClone:  strings.TrimSpace(isShallow) == "true",
	}, nil
}

// Exec supports the execution of any raw git command. No attempt will be
// made to validate the command, and any output will be returned in its
// raw unparsed form
func (c *Client) Exec(cmd string) (string, error) {
	return c.exec(cmd)
}

func (*Client) exec(cmd string) (string, error) {
	p, _ := syntax.NewParser().Parse(strings.NewReader(cmd), "")

	var buf bytes.Buffer
	r, _ := interp.New(
		interp.StdIO(os.Stdin, &buf, &buf),
	)

	if err := r.Run(context.Background(), p); err != nil {
		return "", ErrGitExecCommand{
			Cmd: cmd,
			Out: strings.TrimSuffix(buf.String(), "\n"),
		}
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

func (c *Client) rootDir() (string, error) {
	return c.exec("git rev-parse --show-toplevel")
}

// ToRelativePath determines if a path is relative to the
// working directory of the repository and returns the resolved
// relative path. A [ErrGitNonRelativePath] error will be returned
// if the path exists outside of the working directory.
// [RelativeAtRoot] is returned if the path and working directory
// are equivalent
func (c *Client) ToRelativePath(path string) (string, error) {
	root, err := c.rootDir()
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}

	// Ensure slashes are OS agnostic
	rel = filepath.ToSlash(rel)

	// Reject any paths that are not located within the root repository directory
	if strings.HasPrefix(rel, "../") {
		return "", ErrGitNonRelativePath{
			RootDir:      root,
			TargetPath:   path,
			RelativePath: rel,
		}
	}

	return rel, nil
}
