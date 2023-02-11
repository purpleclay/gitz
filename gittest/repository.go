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

	// grabbed from: https://loremipsum.io/
	fileContent = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
)

// RepositoryOption provides a utility for setting repository options during
// initialization. A repository will always be created with sensible default
// values
type RepositoryOption func(*repositoryOptions)

type repositoryOptions struct {
	Log     []LogEntry
	Files   []file
	Commits []string
}

type file struct {
	Path   string
	Staged bool
}

// WithLog ensures the repository will be initialized with a given snapshot
// of commits and tags. Ideal for initializing a repository with a known
// state. The provided log is parsed using [gittest.ParseLog] and expects
// the log in the following format:
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

// InitRepository will attempt to initialize a test repository capable of
// supporting any git operation. Options can be provided to customize the
// initialization process, changing the default configuration used.
//
// It is important to note, that options will be executed within a
// particular order:
//  1. Log history will be imported
//  2. All local empty commits are made without pushing back to the remote
//  3. All named files will be created and staged if required
//
// Repository creation consists of two phases. First, a bare repository
// is initialized, before being cloned locally. This ensures a fully
// working remote. Without customization, the test repository will
// consist of single commit:
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
	require.NoError(t, os.Chdir(tmpDir))

	Exec(t, fmt.Sprintf("git init --bare --initial-branch %s test.git", DefaultBranch))
	Exec(t, "git clone ./test.git")
	require.NoError(t, os.Chdir("./test"))

	// Ensure default config is set on the repository
	require.NoError(t, setConfig("user.name", DefaultAuthorName))
	require.NoError(t, setConfig("user.email", DefaultAuthorEmail))

	// Initialize the repository so that it is ready for use
	Exec(t, fmt.Sprintf(`git commit --allow-empty -m "%s"`, InitialCommit))
	Exec(t, fmt.Sprintf("git push origin %s", DefaultBranch))

	// Process any provided options to ensure repository is initialized as required
	options := &repositoryOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if len(options.Log) > 0 {
		require.NoError(t, importLog(options.Log))
	}

	for _, commit := range options.Commits {
		Exec(t, fmt.Sprintf(`git commit --allow-empty -m "%s"`, commit))
	}

	for _, f := range options.Files {
		require.NoError(t, tempFile(f.Path, fileContent))
		if f.Staged {
			StageFile(t, f.Path)
		}
	}

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(current))
	})
}

func tempFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}

	return nil
}

func importLog(log []LogEntry) error {
	// It is important to reverse the list as we want to write the log back
	// to the repository using oldest to latest
	for i := len(log) - 1; i >= 0; i-- {
		commitCmd := fmt.Sprintf(`git commit --allow-empty -m "%s"`, log[i].Commit)
		if _, err := exec(commitCmd); err != nil {
			return err
		}

		if log[i].Tag == "" {
			continue
		}

		tagCmd := fmt.Sprintf(`git tag "%s"`, log[i].Tag)
		if _, err := exec(tagCmd); err != nil {
			return err
		}

		pushCmd := fmt.Sprintf(`git push --atomic origin %s "%s"`, DefaultBranch, log[i].Tag)
		if out, err := exec(pushCmd); err != nil {
			fmt.Println(out)
			return err
		}
	}

	return nil
}

func setConfig(key, value string) error {
	configCmd := fmt.Sprintf(`git config %s "%s"`, key, value)
	_, err := exec(configCmd)
	return err
}

// Exec will execute any given git command (expecting no failures) and
// return any received output back to the caller
func Exec(t *testing.T, cmd string) string {
	t.Helper()

	out, err := exec(cmd)
	require.NoError(t, err)

	return out
}

func exec(cmd string) (string, error) {
	p, err := syntax.NewParser().Parse(strings.NewReader(cmd), "")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	r, _ := interp.New(
		interp.StdIO(os.Stdin, &buf, &buf),
	)

	if err := r.Run(context.Background(), p); err != nil {
		return "", errors.New(buf.String())
	}

	return buf.String(), nil
}

// Tags returns a list of all local tags associated with the current
// repository. Raw output is returned from the git command:
//
//	git for-each-ref refs/tags
func Tags(t *testing.T) string {
	t.Helper()
	return Exec(t, "git for-each-ref refs/tags")
}

// RemoteTags returns a list of all tags that have been pushed to the
// remote origin of the current repository. Raw output is returned from
// the git command:
//
//	git ls-remote --tags
func RemoteTags(t *testing.T) string {
	t.Helper()
	return Exec(t, "git ls-remote --tags")
}

// StageFile will attempt to use the provided path to stage a file that
// has been modified. The following git command is executed:
//
//	git add '<path>'
func StageFile(t *testing.T, path string) {
	t.Helper()
	Exec(t, fmt.Sprintf("git add '%s'", path))
}

// Commit a snapshot of all changes within the current repository (working directory)
// without pushing it to the remote. The commit will be associated with the
// provided message. The following git command is executed:
//
//	git commit -m '<message>'
func Commit(t *testing.T, message string) {
	t.Helper()
	Exec(t, fmt.Sprintf("git commit -m '%s'", message))
}

// LastCommit returns the last commit from the git log of the current
// repository. Raw output is returned from the git command:
//
//	git log -n1
func LastCommit(t *testing.T) string {
	t.Helper()
	return Exec(t, "git log -n1")
}

// PorcelainStatus returns a snapshot of the current status of a
// repository (working directory) in an easy to parse format.
// Raw output is returned from the git command:
//
//	git status --porcelain
func PorcelainStatus(t *testing.T) string {
	t.Helper()
	return Exec(t, "git status --porcelain")
}

// LogRemote returns the log history of a repository (working directory)
// as it currently exists on the remote. Any local commit that are not
// pushed, will not appear within this log history. Raw output is
// returned from this command:
//
//	git log --oneline origin/main
func LogRemote(t *testing.T) string {
	t.Helper()
	return Exec(t, fmt.Sprintf("git log --oneline %s", DefaultRemoteBranch))
}

// TagLocal creates a tag that is only tracked locally and will not have
// been pushed back to the remote repository. The following git command
// is executed:
//
//	git tag '<tag>'
func TagLocal(t *testing.T, tag string) {
	t.Helper()
	Exec(t, fmt.Sprintf("git tag '%s'", tag))
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
	return Exec(t, fmt.Sprintf("git show '%s'", object))
}
