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
	"fmt"
	"strings"
)

// LogOption provides a way for setting specific options during a log operation.
// Each supported option can customize the way the log history of the current
// repository (working directory) is processed before retrieval
type LogOption func(*logOptions)

type logOptions struct {
	RefRange string
	LogPaths []string
}

// WithRef provides a starting point other than HEAD (most recent commit)
// when retrieving the log history of the current repository (working
// directory). Typically a reference can be either a commit hash, branch
// name or tag. The output of this option will typically be a shorter,
// fine tuned history. This option is mutually exclusive with
// [git.WithRefRange]. All leading and trailing whitespace are trimmed
// from the reference, allowing empty references to be ignored
func WithRef(ref string) LogOption {
	return func(opts *logOptions) {
		opts.RefRange = strings.TrimSpace(ref)
	}
}

// WithRefRange provides both a start and end point when retrieving a
// focused snapshot of the log history from the current repository
// (working directory). Typically a reference can be either a commit
// hash, branch name or tag. The output of this option will be a shorter,
// fine tuned history, for example, the history between two tags.
// This option is mutually exclusive with [git.WithRef]. All leading
// and trailing whitespace are trimmed from the references, allowing
// empty references to be ignored
func WithRefRange(fromRef string, toRef string) LogOption {
	return func(opts *logOptions) {
		from := strings.TrimSpace(fromRef)
		if from == "" {
			from = "HEAD"
		}

		to := strings.TrimSpace(toRef)
		if to != "" {
			to = fmt.Sprintf("...%s", to)
		}

		opts.RefRange = fmt.Sprintf("%s%s", from, to)
	}
}

// WithPaths allows the log history to be retrieved for any number of
// files and folders within the current repository (working directory).
// Only commits that have had a direct impact on those files and folders
// will be retrieved. Paths to files and folders are relative to the
// root of the repository. All leading and trailing whitespace will be
// trimmed from the file paths, allowing empty paths to be ignored
func WithPaths(paths ...string) LogOption {
	return func(opts *logOptions) {
		opts.LogPaths = make([]string, 0, len(paths))
		opts.LogPaths = append(opts.LogPaths, paths...)
	}
}

// Log retrieves the commit log of the current repository (working directory)
// in an easy to parse format. Options can be provided to customize log
// retrieval, creating a targeted snapshot. By default, the entire history
// from the repository HEAD (most recent commit) will be retrieved. The logs
// are generated using the default git options:
//
//	git log --pretty=oneline --abbrev-commit --no-decorate --no-color
func (c *Client) Log(opts ...LogOption) (string, error) {
	options := &logOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Build command based on the provided options
	var logCmd strings.Builder
	logCmd.WriteString("git log ")

	if options.RefRange != "" {
		logCmd.WriteString(options.RefRange)
	}

	logCmd.WriteString(" --pretty=oneline --abbrev-commit --no-decorate --no-color")

	if len(options.LogPaths) > 0 {
		logCmd.WriteString(" --")
		for _, path := range options.LogPaths {
			logCmd.WriteString(fmt.Sprintf(" '%s'", path))
		}
	}

	return exec(logCmd.String())
}
