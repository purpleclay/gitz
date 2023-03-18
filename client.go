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
	"strings"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
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

// Repository provides a snapshot of the current state of a repository
// (working directory)
type Repository struct {
	ShallowClone  bool
	DetachedHead  bool
	DefaultBranch string
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

	return Repository{
		ShallowClone:  strings.TrimSpace(isShallow) == "true",
		DetachedHead:  strings.TrimSpace(isDetached) == "",
		DefaultBranch: strings.TrimPrefix(defaultBranch, "origin/"),
	}, nil
}

// type Cmd func() Msg
// type Msg interface{}
// type BatchMsg []Cmd
// func Batch(cmds ...Cmd) Cmd {}

// type Cmd func(string) (string, error)
// git.Cmd

// Exec()
// ExecBatch()

// type Cmd func() (string, error)
// func Batch(cmds ...Cmd) Cmd {} // write to the same buffer

// exec(cmd Cmd) (string, error)

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
