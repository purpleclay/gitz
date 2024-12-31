package git

import (
	"strings"
)

// PullOption provides a way for setting specific options while pulling changes
// from the remote. Each supported option can customize how changes are pulled
// from the remote and integrated into the current repository (working directory)
type PullOption func(*pullOptions)

type pullOptions struct {
	Config []string
	fetchOptions
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

// WithFetchTags will fetch all tags from the remote into local tag
// references with the same name
func WithFetchTags() PullOption {
	return func(opts *pullOptions) {
		opts.Tags = true
	}
}

// WithFetchDepthTo will limit the number of commits to be fetched from the
// remotes history. If fetching into a shallow clone of a repository,
// this can be used to shorten or deepen the existing history
func WithFetchDepthTo(depth int) PullOption {
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

// WithFetchIgnoreTags disables local tracking of tags from the remote
func WithFetchIgnoreTags() PullOption {
	return func(opts *pullOptions) {
		opts.NoTags = true
	}
}

// WithPullRefSpecs allows remote references to be cherry-picked and
// fetched into the current repository (working copy) during a pull. A
// reference (or refspec) can be as simple as a name, where git will
// automatically resolve any ambiguity, or as explicit as providing a
// source and destination for reference within the remote. Check out the
// official git documentation on how to write a more complex [refspec]
// [refspec]: https://git-scm.com/docs/git-pull#Documentation/git-pull.txt-ltrefspecgt
func WithPullRefSpecs(refs ...string) PullOption {
	return func(opts *pullOptions) {
		opts.RefSpecs = trim(refs...)
	}
}

// Pull all changes from a remote repository and immediately update the current
// repository (working directory) with those changes. This ensures that your current
// repository keeps track of remote changes and stays in sync
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
	buf.WriteString(options.fetchOptions.String())
	return c.Exec(buf.String())
}
