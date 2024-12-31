package git

import (
	"strconv"
	"strings"
)

// FetchOption provides a way for setting specific options while fetching changes
// from the remote. Each supported option can customize how changes are fetched
// from the remote
type FetchOption func(*fetchOptions)

type fetchOptions struct {
	All       bool
	Config    []string
	Depth     int
	Force     bool
	NoTags    bool
	RefSpecs  []string
	Tags      bool
	Unshallow bool
}

func (o fetchOptions) String() string {
	var buf strings.Builder

	if o.All {
		buf.WriteString(" --all")
	}

	if o.Depth > 0 {
		buf.WriteString(" --depth ")
		buf.WriteString(strconv.Itoa(o.Depth))
	}

	if o.Tags {
		buf.WriteString(" --tags")
	}

	if o.Force {
		buf.WriteString(" --force")
	}

	if o.NoTags {
		buf.WriteString(" --no-tags")
	}

	if o.Unshallow {
		buf.WriteString(" --unshallow")
	}

	if len(o.RefSpecs) > 0 {
		buf.WriteString(" origin ")
		buf.WriteString(strings.Join(o.RefSpecs, " "))
	}

	return buf.String()
}

// WithFetchConfig allows temporary git config to be set while fetching
// changes from the remote. Config set using this approach will override
// any config defined within existing git config files. Config must be
// provided as key value pairs, mismatched config will result in an
// [ErrMissingConfigValue] error. Any invalid paths will result in an
// [ErrInvalidConfigPath] error
func WithFetchConfig(kv ...string) FetchOption {
	return func(opts *fetchOptions) {
		opts.Config = trim(kv...)
	}
}

// WithAll will fetch the latest changes from all tracked remotes
func WithAll() FetchOption {
	return func(opts *fetchOptions) {
		opts.All = true
	}
}

// WithTags will fetch all tags from the remote into local tag
// references with the same name
func WithTags() FetchOption {
	return func(opts *fetchOptions) {
		opts.Tags = true
	}
}

// WithDepthTo will limit the number of commits to be fetched from the
// remotes history to the specified depth. If fetching into a shallow
// clone of a repository, this can be used to shorten or deepen the
// existing history
func WithDepthTo(depth int) FetchOption {
	return func(opts *fetchOptions) {
		opts.Depth = depth
	}
}

// WithForce will force the fetching of a remote branch into a local
// branch with a different name (or refspec). Default behavior within
// git prevents such an operation. Typically used in conjunction with
// the [WithFetchRefSpecs] option
func WithForce() FetchOption {
	return func(opts *fetchOptions) {
		opts.Force = true
	}
}

// WithIgnoreTags disables local tracking of tags from the remote
func WithIgnoreTags() FetchOption {
	return func(opts *fetchOptions) {
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
func WithFetchRefSpecs(refs ...string) FetchOption {
	return func(opts *fetchOptions) {
		opts.RefSpecs = trim(refs...)
	}
}

// WithUnshallow will fetch the complete history from the remote
func WithUnshallow() FetchOption {
	return func(opts *fetchOptions) {
		opts.Unshallow = true
	}
}

// Fetch all remote changes from a remote repository without integrating (merging)
// them into the current repository (working directory). Ensures the current repository
// only tracks the latest remote changes
func (c *Client) Fetch(opts ...FetchOption) (string, error) {
	options := &fetchOptions{}
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

	buf.WriteString(" fetch")
	buf.WriteString(options.String())
	return c.Exec(buf.String())
}
