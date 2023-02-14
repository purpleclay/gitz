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
type TagOption func(*tagOptions)

type tagOptions struct {
	Annotation string
}

// WithAnnotation ensures the created tag is annotated with the provided
// message. This ultimately converts the standard lightweight tag into
// an annotated tag which is stored as a full object within the git
// database. Any leading and trailing whitespace will automatically be
// trimmed from the message. This allows empty messages to be ignored
func WithAnnotation(message string) TagOption {
	return func(opts *tagOptions) {
		opts.Annotation = strings.TrimSpace(message)
	}
}

// Tag a specific point within a repositories history and push it to the
// configured remote. Tagging comes in two flavours:
//   - A lightweight tag, which points to a specific commit within
//     the history and marks a specific point in time
//   - An annotated tag, which is treated as a full object within
//     git, and must include a tagging message (or annotation)
//
// By default, a lightweight tag will be created, unless specific tag
// options are provided
func (c *Client) Tag(tag string, opts ...TagOption) (string, error) {
	options := &tagOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Build command based on the provided options
	var tagCmd strings.Builder
	tagCmd.WriteString("git tag ")

	if options.Annotation != "" {
		tagCmd.WriteString(fmt.Sprintf("-m '%s' -a ", options.Annotation))
	}
	tagCmd.WriteString(fmt.Sprintf("'%s'", tag))

	if out, err := exec(tagCmd.String()); err != nil {
		return out, err
	}

	return exec(fmt.Sprintf("git push origin '%s'", tag))
}

// DeleteTag a tag both locally and from the remote origin
func (c *Client) DeleteTag(tag string) (string, error) {
	if out, err := exec(fmt.Sprintf(`git tag -d "%s"`, tag)); err != nil {
		return out, err
	}

	return exec(fmt.Sprintf("git push --delete origin '%s'", tag))
}
