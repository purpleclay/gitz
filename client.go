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
	ShallowClone  bool
	DetachedHead  bool
	DefaultBranch string
	RootDir       string
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
	if _, err := exec("type git"); err != nil {
		return nil, ErrGitMissing{PathEnv: os.Getenv("PATH")}
	}

	c := &Client{}
	c.gitVersion, _ = exec("git --version")
	return c, nil
}

// Version of git used by the client
func (c *Client) Version() string {
	return c.gitVersion
}

// Repository returns details about the current repository (working directory),
// by carrying out a series of checks. Answers to which are returned as a
// snapshot for querying
func (c *Client) Repository() (Repository, error) {
	isRepo, _ := exec("git rev-parse --is-inside-work-tree")
	if strings.TrimSpace(isRepo) != "true" {
		return Repository{}, errors.New("current working directory is not a git repository")
	}

	isShallow, _ := exec("git rev-parse --is-shallow-repository")
	isDetached, _ := exec("git branch --show-current")
	defaultBranch, _ := exec("git rev-parse --abbrev-ref remotes/origin/HEAD")
	rootDir, _ := c.rootDir()

	return Repository{
		ShallowClone:  strings.TrimSpace(isShallow) == "true",
		DetachedHead:  strings.TrimSpace(isDetached) == "",
		DefaultBranch: strings.TrimPrefix(defaultBranch, "origin/"),
		RootDir:       rootDir,
	}, nil
}

func exec(cmd string) (string, error) {
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
	return exec("git rev-parse --show-toplevel")
}

// ToRelativePath determines if a path is relative to the
// working directory of the repository and returns the resolved
// relative path. A [ErrGitNonRelativePath] error will be returned
// if the path exists outside of the working directory
func (c *Client) ToRelativePath(path string) (string, error) {
	root, err := c.rootDir()
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}

	// Reject any paths that are not located within the root repository directory
	if strings.HasPrefix(filepath.ToSlash(rel), "../") {
		return "", ErrGitNonRelativePath{
			RootDir:      root,
			TargetPath:   path,
			RelativePath: rel,
		}
	}

	return rel, nil
}
