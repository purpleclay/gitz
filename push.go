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

// PushOption provides a way of setting specific options during a git
// push operation. Each supported option can customize the way in which
// references are pushed back to the remote
type PushOption func(*pushOptions)

type pushOptions struct {
	All         bool
	PushOptions []string
	Tags        bool
	RefSpecs    []string
}

// WithAllBranches will push all locally created branch references
// back to the remote
func WithAllBranches() PushOption {
	return func(opts *pushOptions) {
		opts.All = true
	}
}

// WithAllTags will push all locally created tag references back
// to the remote
func WithAllTags() PushOption {
	return func(opts *pushOptions) {
		opts.Tags = true
	}
}

// WithPushOptions allows any number of aribitrary strings to be pushed
// to the remote server. All options are transmitted in their received
// order. A server must have the git config setting receive.advertisePushOptions
// set to true to receive push options
func WithPushOptions(options ...string) PushOption {
	return func(opts *pushOptions) {
		opts.PushOptions = trim(options...)
	}
}

// WithRefSpecs allows locally created references to be cherry-picked
// and pushed back to the remote. A reference (or refspec) can be as
// simple as a name, where git will automatically resolve any
// ambiguity, or as explicit as providing a source and destination
// for each local reference within the remote. Check out the official
// git documentation on how to write a more complex [refspec]
//
// [refspec]: https://git-scm.com/docs/git-push#Documentation/git-push.txt-ltrefspecgt82308203
func WithRefSpecs(refs ...string) PushOption {
	return func(opts *pushOptions) {
		opts.RefSpecs = trim(refs...)
	}
}

// Push (or upload) all local changes to the remote repository.
// By default, changes associated with the current branch will
// be pushed back to the remote. Options can be provided to
// configure branch and tag push semantics
func (c *Client) Push(opts ...PushOption) (string, error) {
	options := &pushOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var buffer strings.Builder
	buffer.WriteString("git push")

	for _, po := range options.PushOptions {
		buffer.WriteString(" --push-option=" + po)
	}

	if options.All {
		buffer.WriteString(" --all")
	} else if options.Tags {
		buffer.WriteString(" --tags")
	} else if len(options.RefSpecs) > 0 {
		buffer.WriteString(" origin ")
		buffer.WriteString(strings.Join(options.RefSpecs, " "))
	} else {
		out, err := c.exec("git branch --show-current")
		if err != nil {
			return out, err
		}

		buffer.WriteString(fmt.Sprintf(" origin %s", out))
	}

	return c.exec(buffer.String())
}

// PushRefOption provides a way of setting specific options during a
// git push of a specific reference. Each supported option can customize
// the way in which a reference is pushed back to the remote
type PushRefOption func(*pushRefOptions)

type pushRefOptions struct {
	Delete bool
}

// WithRefDelete will trigger the deletion of a reference when pushed
// back to the remote
func WithRefDelete() PushRefOption {
	return func(opts *pushRefOptions) {
		opts.Delete = true
	}
}

// PushRef will push an individual reference to the remote repository
func (c *Client) PushRef(ref string, opts ...PushRefOption) (string, error) {
	return c.PushRefs([]string{ref}, opts...)
}

// PushRefs will push a batch of references to the remote repository
func (c *Client) PushRefs(refs []string, opts ...PushRefOption) (string, error) {
	if len(refs) == 0 {
		return "", nil
	}

	options := &pushRefOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var buffer strings.Builder
	buffer.WriteString("git push origin ")

	if options.Delete {
		buffer.WriteString("--delete ")
	}

	buffer.WriteString(strings.Join(refs, " "))
	return c.exec(buffer.String())
}
