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

// PullOption provides a way for setting specific options while pulling changes
// from the remote. Each supported option can customize how changes are pulled
// from the remote and integrated into the current repository (working directory)
type PullOption func(*pullOptions)

type pullOptions struct {
	All      bool
	AllTags  bool
	Config   []string
	Depth    int
	Force    bool
	NoTags   bool
	RefSpecs []string
}

// WithPullConfig allows temporary git config to be set while pulling
// changes from the remote. Config set using this approach will override
// any config defined within existing git config files. Config must be
// provided as key value pairs, mismatched config will result in an
// [ErrMissingConfigValue] error. Any invalid paths will result in an
// [ErrInvalidConfigPath] error
func WithPullConfig(kv ...string) PullOption {
	return func(opts *pullOptions) {
		opts.Config = trim(kv...)
	}
}

// WithFetchAll will fetch the latest changes from all tracked remotes
func WithFetchAll() PullOption {
	return func(opts *pullOptions) {
		opts.All = true
	}
}

// WithFetchAllTags will fetch all tags from the remote into local tag
// references with the same name
func WithFetchAllTags() PullOption {
	return func(opts *pullOptions) {
		opts.AllTags = true
	}
}

// WithFetchDepth will limit the number of commits to be fetched from the
// remotes history. If fetching into a shallow clone of a repository,
// this can be used to shorten or deepen the existing history
func WithFetchDepth(depth int) PullOption {
	return func(opts *pullOptions) {
		opts.Depth = depth
	}
}

// WithFetchForce will force the fetching of a remote branch into a local
// branch with a different name (or refspec). Default behavior within
// git prevents such an operation. Typically used in conjunction with
// the [WithFetchRefSpecs] option
func WithFetchForce() PullOption {
	return func(opts *pullOptions) {
		opts.Force = true
	}
}

// WithFetchNoTags disables local tracking of tags from the remote
func WithFetchNoTags() PullOption {
	return func(opts *pullOptions) {
		opts.NoTags = true
	}
}

// WithFetchRefSpecs allows remote references to be cherry-picked and
// fetched into the current repository (working copy). A reference
// (or refspec) can be as simple as a name, where git will automatically
// resolve any ambiguity, or as explicit as providing a source and destination
// for reference within the remote. Check out the official git documentation
// on how to write a more complex [refspec]
// [refspec]: https://git-scm.com/docs/git-fetch#Documentation/git-fetch.txt-ltrefspecgt
func WithFetchRefSpecs(refs ...string) PullOption {
	return func(opts *pullOptions) {
		opts.RefSpecs = trim(refs...)
	}
}

// Pull all changes from a remote repository and immediately update the current
// repository (current working) directory with those changes. This ensures
// that your current repository keeps track of remote changes and stays in sync
func (c *Client) Pull(opts ...PullOption) (string, error) {
	options := &pullOptions{}
	for _, opt := range opts {
		opt(options)
	}

	cfg, err := ToInlineConfig(options.Config...)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	buf.WriteString("git")

	if len(cfg) > 0 {
		buf.WriteString(" ")
		buf.WriteString(strings.Join(cfg, " "))
	}
	buf.WriteString(" pull")

	if options.All {
		buf.WriteString(" --all")
	}

	if options.Depth > 0 {
		buf.WriteString(" --depth ")
		buf.WriteString(strconv.Itoa(options.Depth))
	}

	if options.AllTags {
		buf.WriteString(" --tags")
	}

	if options.Force {
		buf.WriteString(" --force")
	}

	if options.NoTags {
		buf.WriteString(" --no-tags")
	}

	if len(options.RefSpecs) > 0 {
		buf.WriteString(" origin ")
		buf.WriteString(strings.Join(options.RefSpecs, " "))
	}

	return c.exec(buf.String())
}
