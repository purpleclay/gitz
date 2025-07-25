package git

import (
	"fmt"
	"strings"
)

// PushOption provides a way of setting specific options during a git
// push operation. Each supported option can customize the way in which
// references are pushed back to the remote.
type PushOption func(*pushOptions)

type pushOptions struct {
	All         bool
	Config      []string
	Delete      bool
	PushOptions []string
	Tags        bool
	RefSpecs    []string
}

// WithAllBranches will push all locally created branch references
// back to the remote.
func WithAllBranches() PushOption {
	return func(opts *pushOptions) {
		opts.All = true
	}
}

// WithAllTags will push all locally created tag references back
// to the remote.
func WithAllTags() PushOption {
	return func(opts *pushOptions) {
		opts.Tags = true
	}
}

// WithDeleteRefSpecs will trigger the deletion of any named references
// when pushed back to the remote.
func WithDeleteRefSpecs(refs ...string) PushOption {
	return func(opts *pushOptions) {
		opts.Delete = true
		opts.RefSpecs = trim(refs...)
	}
}

// WithPushConfig allows temporary git config to be set while pushing
// changes to the remote. Config set using this approach will override
// any config defined within existing git config files. Config must be
// provided as key value pairs, mismatched config will result in an
// [ErrMissingConfigValue] error. Any invalid paths will result in an
// [ErrInvalidConfigPath] error.
func WithPushConfig(kv ...string) PushOption {
	return func(opts *pushOptions) {
		opts.Config = trim(kv...)
	}
}

// WithPushOptions allows any number of aribitrary strings to be pushed
// to the remote server. All options are transmitted in their received
// order. A server must have the git config setting receive.advertisePushOptions
// set to true to receive push options.
func WithPushOptions(options ...string) PushOption {
	return func(opts *pushOptions) {
		opts.PushOptions = trim(options...)
	}
}

// WithRefSpecs allows local references to be cherry-picked and
// pushed back to the remote. A reference (or refspec) can be as
// simple as a name, where git will automatically resolve any
// ambiguity, or as explicit as providing a source and destination
// for each local reference within the remote. Check out the official
// git documentation on how to write a more complex [refspec].
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
// configure branch and tag push semantics.
func (c *Client) Push(opts ...PushOption) (string, error) {
	options := &pushOptions{}
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
	buf.WriteString(" push")

	for _, po := range options.PushOptions {
		buf.WriteString(" --push-option=" + po)
	}

	//nolint:gocritic
	if options.All {
		buf.WriteString(" --all")
	} else if options.Tags {
		buf.WriteString(" --tags")
	} else if len(options.RefSpecs) > 0 {
		buf.WriteString(" origin ")
		if options.Delete {
			buf.WriteString("--delete ")
		}

		buf.WriteString(strings.Join(options.RefSpecs, " "))
	} else {
		out, err := c.Exec("git branch --show-current")
		if err != nil {
			return out, err
		}
		buf.WriteString(fmt.Sprintf(" origin %s", out))
	}

	return c.Exec(buf.String())
}

// PushRef will push an individual reference to the remote repository
// Deprecated: use [Push] instead.
func (c *Client) PushRef(ref string) (string, error) {
	return c.Exec(fmt.Sprintf("git push origin %s", ref))
}
