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

// PushOption ...
type PushOption func(*pushOptions)

type pushOptions struct {
	All      bool
	Tags     bool
	RefSpecs []string
}

// WithAllBranches ...
func WithAllBranches() PushOption {
	return func(opts *pushOptions) {
		opts.All = true
	}
}

// WithAllTags ...
func WithAllTags() PushOption {
	return func(opts *pushOptions) {
		opts.Tags = true
	}
}

// WithRefSpecs ...
func WithRefSpecs(refs ...string) PushOption {
	return func(opts *pushOptions) {
		opts.RefSpecs = Trim(refs...)
	}
}

// Push (or upload) all local changes to the remote repository
func (c *Client) Push(opts ...PushOption) (string, error) {
	options := &pushOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var buffer strings.Builder
	buffer.WriteString("git push")

	if options.All {
		buffer.WriteString(" --all")
	} else if options.Tags {
		buffer.WriteString(" --tags")
	} else if len(options.RefSpecs) > 0 {
		buffer.WriteString(" origin ")
		buffer.WriteString(strings.Join(options.RefSpecs, " "))
	} else {
		out, err := exec("git branch --show-current")
		if err != nil {
			return out, err
		}

		buffer.WriteString(fmt.Sprintf(" origin '%s'", out))
	}

	return exec(buffer.String())
}

// PushTag will push an individual tag reference to the remote repository
func (c *Client) PushTag(tag string) (string, error) {
	return exec(fmt.Sprintf("git push origin '%s'", tag))
}
