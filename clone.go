package git

import (
	"strconv"
	"strings"
)

// CloneOption provides a way for setting specific options during a clone
// operation. Each supported option can customize the way in which the
// repository is cloned onto the file system into a target working directory
type CloneOption func(*cloneOptions)

type cloneOptions struct {
	Config      []string
	CheckoutRef string
	Depth       int
	Dir         string
	NoTags      bool
}

// WithCheckoutRef changes the default checkout behavior after a clone succeeds.
// A branch or tag reference is supported. Checking out a tag will result in
// a detached HEAD. An empty string will be ignored
func WithCheckoutRef(ref string) CloneOption {
	return func(opts *cloneOptions) {
		opts.CheckoutRef = strings.TrimSpace(ref)
	}
}

// WithCloneConfig allows temporary git config to be set while cloning
// the remote into a newly created directory. Config set using this
// approach will override any config defined within existing git config
// files. Config must be provided as key value pairs, mismatched config
// will result in an [ErrMissingConfigValue] error. Any invalid paths will
// result in an [ErrInvalidConfigPath] error
func WithCloneConfig(kv ...string) CloneOption {
	return func(opts *cloneOptions) {
		opts.Config = trim(kv...)
	}
}

// WithDepth ensures the repository will be cloned at a specific depth,
// effectively truncating the history to the required number of commits.
// The result will be a shallow repository. Any depth less than one
// is ignored, resulting in a full clone of the history
func WithDepth(depth int) CloneOption {
	return func(opts *cloneOptions) {
		opts.Depth = depth
	}
}

// WithDirectory provides a named directory for cloning the repository into.
// If the directory already exists, it must be empty for the clone to
// be successful. An empty string will be ignored
func WithDirectory(dir string) CloneOption {
	return func(opts *cloneOptions) {
		opts.Dir = strings.TrimSpace(dir)
	}
}

// WithNoTags prevents any tags from being included during the clone
func WithNoTags() CloneOption {
	return func(opts *cloneOptions) {
		opts.NoTags = true
	}
}

// Clone a repository by its provided URL into a newly created directory.
// A default clone will ensure remote tracking branches are created for
// each branch within the repository with only the default branch being
// checked out fully. The cloned directory will mirror that of the repository
// name within its URL. Options can be provided to customize the clone
// behavior
func (c *Client) Clone(url string, opts ...CloneOption) (string, error) {
	options := &cloneOptions{}
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
	buf.WriteString(" clone")

	if options.NoTags {
		buf.WriteString(" --no-tags")
	}

	if options.CheckoutRef != "" {
		buf.WriteString(" --branch ")
		buf.WriteString(options.CheckoutRef)
	}

	if options.Depth > 0 {
		buf.WriteString(" --depth ")
		buf.WriteString(strconv.Itoa(options.Depth))
	}

	buf.WriteString(" -- ")
	buf.WriteString(url)

	if options.Dir != "" {
		buf.WriteRune(' ')
		buf.WriteString(options.Dir)
	}

	return c.Exec(buf.String())
}
