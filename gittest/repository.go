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

package gittest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

const (
	// DefaultBranch contains the name of the default branch used when
	// initializing the test repository
	DefaultBranch = "main"

	// DefaultRemoteBranch contains the name of the default branch when
	// initializing the remote bare repository
	DefaultRemoteBranch = "origin/main"

	// DefaultAuthorName contains the author name written to local git
	// config when initializing the test repository
	DefaultAuthorName = "batman"

	// DefaultAuthorEmail contains the author email written to local git
	// config when initializing the test repository
	DefaultAuthorEmail = "batman@dc.com"

	// DefaultAuthorLog contains the default git representation of an author
	// and can be used for matching against entries within a git log
	DefaultAuthorLog = "batman <batman@dc.com>"

	// InitialCommit contains the first commit message used to initialize
	// the test repository
	InitialCommit = "initialized repository"

	// BareRepositoryName the name of the bare repository, used as the
	// remote for all testing
	BareRepositoryName = "test.git"

	// ClonedRepositoryName the name of the repository (working directory)
	// after cloning the bare repository
	ClonedRepositoryName = "test"

	// grabbed from: https://loremipsum.io/
	fileContent = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	// an internal template for pushing changes back to a remote origin
	gitPushTemplate = "git push origin %s"
)

// RepositoryOption provides a utility for setting repository options during
// initialization. A repository will always be created with sensible default
// values
type RepositoryOption func(*repositoryOptions)

type repositoryOptions struct {
	Log        []LogEntry
	RemoteLog  []LogEntry
	Files      []file
	Commits    []string
	CloneDepth int
}

type file struct {
	Path   string
	Staged bool
}

// WithLog ensures the repository will be initialized with a given snapshot
// of commits and tags. Ideal for initializing a repository with a known
// state. The provided log is parsed using [ParseLog] and expects the log
// in the following format:
//
//	(tag: 0.1.0) feat: improve existing cli documentation
//	docs: create initial mkdocs material documentation
//
// This is the equivalent to the format produced using the git command:
//
//	git log --pretty='format:%d %s'
func WithLog(log string) RepositoryOption {
	return func(opts *repositoryOptions) {
		opts.Log = ParseLog(log)
	}
}

// WithRemoteLog ensures the remote origin of the repository will be
// initialized with a given snapshot of commits and tags. Ideal for
// simulating a delta between the current repository (working directory)
// and the remote. Use with caution, as this can result in conflicts.
// The provided log is parsed using [ParseLog] and expects the log
// in the following format:
//
//	(tag: 0.1.0) feat: improve existing cli documentation
//	docs: create initial mkdocs material documentation
//
// This is the equivalent to the format produced using the git command:
//
//	git log --pretty='format:%d %s'
func WithRemoteLog(log string) RepositoryOption {
	return func(opts *repositoryOptions) {
		opts.RemoteLog = ParseLog(log)
	}
}

// WithFiles ensures the repository will be initialized with a given set
// of named files. Both relative and full file paths are supported. Each
// file will be generated using default data, but will remain untracked
// by the repository.
//
// For example:
//
//	gittest.InitRepository(t, gittest.WithFiles("file1.txt", "file2.txt"))
//
// This will result in a repository containing two untracked files. Which
// can be verified by checking the git status:
//
//	$ git status --porcelain
//	?? file1.txt
//	?? file2.txt
func WithFiles(files ...string) RepositoryOption {
	return func(opts *repositoryOptions) {
		for _, f := range files {
			opts.Files = append(opts.Files, file{Path: f, Staged: false})
		}
	}
}

// WithStagedFiles ensures the repository will be initialized with a given
// set of named files. Both relative and full file paths are supported. Each
// file will be generated using default data, and will be staged within the
// repository.
//
// For example:
//
//	gittest.InitRepository(t, gittest.WithStagedFiles("file1.txt", "file2.txt"))
//
// This will result in a repository containing two staged files. Which
// can be verified by checking the git status:
//
//	$ git status --porcelain
//	A  file1.txt
//	A  file2.txt
func WithStagedFiles(files ...string) RepositoryOption {
	return func(opts *repositoryOptions) {
		for _, f := range files {
			opts.Files = append(opts.Files, file{Path: f, Staged: true})
		}
	}
}

// WithLocalCommits ensures the repository will be initialized with a set
// of local empty commits, which will not have been pushed back to the remote
func WithLocalCommits(commits ...string) RepositoryOption {
	return func(opts *repositoryOptions) {
		opts.Commits = commits
	}
}

// WithCloneDepth ensures the repository will be cloned at a specific depth,
// effectively truncating the history to the required number of commits.
// The result will be a shallow repository
func WithCloneDepth(depth int) RepositoryOption {
	return func(opts *repositoryOptions) {
		opts.CloneDepth = depth
	}
}

// InitRepository will attempt to initialize a test repository capable of
// supporting any git operation. Options can be provided to customize the
// initialization process, changing the default configuration used.
//
// It is important to note, that options will be executed within a
// particular order:
//  1. Log history will be imported (local and remote are in sync)
//  2. A shallow clone is made at the required clone depth
//  3. Remote log history will be imported, creating a delta between
//     the current repository (working directory) and the remote
//  4. All local empty commits are made without pushing back to the remote
//  5. All named files will be created and staged if required
//
// Repository creation consists of two phases. First, a bare repository
// is initialized, before being cloned locally. This ensures a fully
// working remote. Without customization (options), the test repository
// will consist of single commit:
//
//	initialized repository
func InitRepository(t *testing.T, opts ...RepositoryOption) {
	t.Helper()

	// Track our current directory
	current, err := os.Getwd()
	require.NoError(t, err)

	// Generate two temporary directories. The first is initialized as a
	// bare repository and becomes our filesystem based remote. The second
	// is our working repository, which is a clone of the former
	tmpDir := t.TempDir()
	changeToDir(t, tmpDir)

	Exec(t, fmt.Sprintf("git init --bare --initial-branch %s %s", DefaultBranch, BareRepositoryName))
	cloneRemoteAndInit(t, ClonedRepositoryName)

	// Process any provided options to ensure repository is initialized as required
	options := &repositoryOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if len(options.Log) > 0 {
		importLog(t, options.Log)
	}

	if options.CloneDepth > 0 {
		// Remove the existing local clone and clone again specifying the depth
		changeToDir(t, tmpDir)
		require.NoError(t, os.RemoveAll(ClonedRepositoryName))
		cloneRemoteAndInit(t, ClonedRepositoryName, fmt.Sprintf("--depth %d", options.CloneDepth))
	}

	// To ensure a successful delta is created, an additional clone is made of the
	// bare (remote) repository. The remote log is then imported, ensuring the
	// local clone is out of sync
	if len(options.RemoteLog) > 0 {
		localClone := changeToDir(t, tmpDir)
		cloneRemoteAndInit(t, "remote-import")

		importLog(t, options.RemoteLog)
		require.NoError(t, os.Chdir(localClone))
	}

	for _, commit := range options.Commits {
		Exec(t, fmt.Sprintf(`git commit --allow-empty -m "%s"`, commit))
	}

	for _, f := range options.Files {
		tempFile(t, f.Path, fileContent)
		if f.Staged {
			StageFile(t, f.Path)
		}
	}

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(current))
	})
}

func changeToDir(t *testing.T, dir string) string {
	changedFrom, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(dir))
	return changedFrom
}

func cloneRemoteAndInit(t *testing.T, cloneName string, options ...string) {
	MustExec(t, fmt.Sprintf("git clone %s file://$(pwd)/%s %s", strings.Join(options, " "), BareRepositoryName, cloneName))
	require.NoError(t, os.Chdir(cloneName))

	// Ensure author details are set
	require.NoError(t, setConfig("user.name", DefaultAuthorName))
	require.NoError(t, setConfig("user.email", DefaultAuthorEmail))

	// Check if there any any commits, if not, initialize and push back first commit
	if out := MustExec(t, "git rev-list -n1 --all"); out == "" {
		MustExec(t, fmt.Sprintf(`git commit --allow-empty -m "%s"`, InitialCommit))
		MustExec(t, fmt.Sprintf(gitPushTemplate, DefaultBranch))
	}

	MustExec(t, "git remote set-head origin --auto")
}

func tempFile(t *testing.T, path, content string) {
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func importLog(t *testing.T, log []LogEntry) {
	// It is important to reverse the list as we want to write the log back
	// to the repository using oldest to latest
	for i := len(log) - 1; i >= 0; i-- {
		commitCmd := fmt.Sprintf(`git commit --allow-empty -m "%s"`, log[i].Commit)
		MustExec(t, commitCmd)

		commitPushCmd := fmt.Sprintf(gitPushTemplate, DefaultBranch)
		MustExec(t, commitPushCmd)

		// If there is no tag, continue onto the next log entry
		if log[i].Tag == "" {
			continue
		}

		tagCmd := fmt.Sprintf("git tag %s", log[i].Tag)
		MustExec(t, tagCmd)

		pushTagCmd := fmt.Sprintf(gitPushTemplate, log[i].Tag)
		MustExec(t, pushTagCmd)
	}
}

func setConfig(key, value string) error {
	configCmd := fmt.Sprintf(`git config %s "%s"`, key, value)
	_, err := exec(configCmd)
	return err
}

// Exec will execute any given git command and return the raw output and
// error from the underlying git client
func Exec(t *testing.T, cmd string) (string, error) {
	t.Helper()
	return exec(cmd)
}

// MustExec will execute any given git command, requiring no failure. Any
// raw output will be returned from the underlying git client
func MustExec(t *testing.T, cmd string) string {
	t.Helper()

	out, err := exec(cmd)
	require.NoError(t, err)

	return out
}

func exec(cmd string) (string, error) {
	p, _ := syntax.NewParser().Parse(strings.NewReader(cmd), "")

	var buf bytes.Buffer
	r, _ := interp.New(
		interp.StdIO(os.Stdin, &buf, &buf),
	)

	if err := r.Run(context.Background(), p); err != nil {
		return "", errors.New(strings.TrimSuffix(buf.String(), "\n"))
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// Tags returns a list of all local tags associated with the current
// repository. Raw output is returned from the git command:
//
//	git for-each-ref refs/tags
func Tags(t *testing.T) string {
	t.Helper()
	return MustExec(t, "git for-each-ref refs/tags")
}

// RemoteTags returns a list of all tags that have been pushed to the
// remote origin of the current repository. Raw output is returned from
// the git command:
//
//	git ls-remote --tags
func RemoteTags(t *testing.T) string {
	t.Helper()
	return MustExec(t, "git ls-remote --tags")
}

// StageFile will attempt to use the provided path to stage a file that
// has been modified. The following git command is executed:
//
//	git add '<path>'
func StageFile(t *testing.T, path string) {
	t.Helper()
	MustExec(t, fmt.Sprintf("git add '%s'", path))
}

// Commit a snapshot of all changes within the current repository (working directory)
// without pushing it to the remote. The commit will be associated with the
// provided message. The following git command is executed:
//
//	git commit -m '<message>'
func Commit(t *testing.T, message string) {
	t.Helper()
	MustExec(t, fmt.Sprintf("git commit -m '%s'", message))
}

// LastCommit returns the last commit from the git log of the current
// repository. Raw output is returned from the git command:
//
//	git log -n1
func LastCommit(t *testing.T) string {
	t.Helper()
	return MustExec(t, "git log -n1")
}

// PorcelainStatus returns a snapshot of the current status of a
// repository (working directory) in an easy to parse format.
// Raw output is returned from the git command:
//
//	git status --porcelain
func PorcelainStatus(t *testing.T) string {
	t.Helper()
	return MustExec(t, "git status --porcelain")
}

// LogRemote returns the log history of a repository (working directory)
// as it currently exists on the remote. Any local commit that are not
// pushed, will not appear within this log history. Raw output is
// returned from this command:
//
//	git log --oneline origin/main
func LogRemote(t *testing.T) string {
	t.Helper()
	return MustExec(t, fmt.Sprintf("git log --oneline %s", DefaultRemoteBranch))
}

// TagLocal creates a tag that is only tracked locally and will not have
// been pushed back to the remote repository. The following git command
// is executed:
//
//	git tag '<tag>'
func TagLocal(t *testing.T, tag string) {
	t.Helper()
	MustExec(t, fmt.Sprintf("git tag '%s'", tag))
}

// Show will display information about a specific git object. The output
// will vary based on the type of object being shown:
//   - For commits it shows the log message and textual diff
//   - For tags, it shows the tag message and the referenced objects
//   - For trees, it shows the names
//   - For plain blobs, it shows the plain contents
//
// Raw output is returned from this command:
//
//	git show '<object>'
func Show(t *testing.T, object string) string {
	t.Helper()
	return MustExec(t, fmt.Sprintf("git show '%s'", object))
}

// Checkout will update the state of the repository (working directory)
// by updating files in the tree to a specific point in time. Object
// can be any one of the following:
//   - Commit reference (long or abbreviated hash)
//   - Tag reference
//   - Branch name
//
// Be warned, checking out a tag or commit reference will cause a
// detached HEAD for the current repository. Raw output is returned
// from this command:
//
//	git checkout '<object>'
func Checkout(t *testing.T, object string) string {
	t.Helper()
	return MustExec(t, fmt.Sprintf("git checkout '%s'", object))
}

// Remote will retrieve the URL of the remote (typically origin) configured
// for the current repository (working directory). To prevent issues due
// to OS dependent separators, the raw URL will be converted to use the
// '/' separator which is compatible across OS when using the git client.
//
// Remote is queried using this command:
//
//	git ls-remote --get-url
func Remote(t *testing.T) string {
	t.Helper()
	remote := MustExec(t, "git ls-remote --get-url")

	// Ensure path is escaped correctly when testing across different OS
	return filepath.ToSlash(remote)
}
