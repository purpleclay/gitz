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

// TagOption provides a way for setting specific options during a tag operation.
// Each supported option can customize the way the tag is applied against
// the current repository (working directory)

// LogOption ...
type LogOption func(*logOptions)

type logOptions struct {
	RefRange string
}

// TODO: Mutually exclusive WithRef and WithRefRange overwrite each other

// WithRef ...
func WithRef(ref string) LogOption {
	return func(opts *logOptions) {
		opts.RefRange = strings.TrimSpace(ref)
	}
}

// WithRefRange ...
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

// TODO: rewrite function docs

// Log retrieves the commit log of the current repository (working directory)
// in an easy to parse format. The logs are generated using the default
// git options:
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

	return exec(logCmd.String())
}
