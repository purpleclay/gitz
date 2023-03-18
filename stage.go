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

import "strings"

// StageOption provides a way for setting specific options during a stage
// operation. Each supported option can customize the way files are staged
// within the current repository (working directory)
type StageOption func(*stageOptions)

type stageOptions struct {
	PathSpecs []string
}

// WithPathSpecs permits a series of [PathSpecs] (or globs) to be defined
// that will stage any matching files within the current repository
// (working directory). Paths to files and folders are relative to the
// root of the repository. All leading and trailing whitespace will be
// trimmed from the file paths, allowing empty paths to be ignored
//
// [PathSpecs]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddefpathspecapathspec
func WithPathSpecs(specs ...string) StageOption {
	return func(opts *stageOptions) {
		opts.PathSpecs = Trim(specs...)
	}
}

// Stage changes to any file or folder within the current repository
// (working directory) ready for inclusion in the next commit. Options
// can be provided to further configure stage semantics. By default,
// all changes will be staged ready for the next commit.
func (c *Client) Stage(opts ...StageOption) (string, error) {
	options := &stageOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Build command based on the provided options
	var stageCmd strings.Builder
	stageCmd.WriteString("git add ")

	if len(options.PathSpecs) > 0 {
		stageCmd.WriteString("--")
		for _, spec := range options.PathSpecs {
			stageCmd.WriteString(" ")
			stageCmd.WriteString(spec)
		}
	} else {
		stageCmd.WriteString("--all")
	}

	return exec(stageCmd.String())
}
