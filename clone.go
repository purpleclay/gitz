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
	"strconv"
	"strings"
)

// CloneOption ...
type CloneOption func(*cloneOptions)

type cloneOptions struct {
	Branch string
	Depth  int
	Dir    string
	NoTags bool
}

// WithBranch ...
func WithBranch(branch string) CloneOption {
	return func(opts *cloneOptions) {
		opts.Branch = strings.TrimSpace(branch)
	}
}

// WithDepth ...
func WithDepth(depth int) CloneOption {
	return func(opts *cloneOptions) {
		opts.Depth = depth
	}
}

// WithDirectory ...
func WithDirectory(dir string) CloneOption {
	return func(opts *cloneOptions) {
		opts.Dir = strings.TrimSpace(dir)
	}
}

// WithNoTags ...
func WithNoTags() CloneOption {
	return func(opts *cloneOptions) {
		opts.NoTags = true
	}
}

// Clone a repository by its provided URL into a newly created directory.
// Remote tracking branches are created for each branch within the repository
// with only the default branch being checked out fully
func (c *Client) Clone(url string, opts ...CloneOption) (string, error) {
	options := &cloneOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var buffer strings.Builder
	buffer.WriteString("git clone")
	if options.NoTags {
		buffer.WriteString(" --no-tags")
	}

	if options.Branch != "" {
		buffer.WriteString(" --branch ")
		buffer.WriteString(options.Branch)
	}

	if options.Depth > 0 {
		buffer.WriteString(" --depth ")
		buffer.WriteString(strconv.Itoa(options.Depth))
	}

	buffer.WriteString(" -- ")
	buffer.WriteString(url)

	if options.Dir != "" {
		buffer.WriteRune(' ')
		buffer.WriteString(options.Dir)
	}

	return exec(buffer.String())
}
