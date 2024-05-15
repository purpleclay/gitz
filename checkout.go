package git

import (
	"strings"
)

// CheckoutOption provides a way for setting specific options while attempting
// to checkout a branch. Each supported option can customize how a branch is checked
// out from the remote and integrated into the current repository (working directory)
type CheckoutOption func(*checkoutOptions)

type checkoutOptions struct {
	Config []string
}

// WithCheckoutConfig allows temporary git config to be set while checking
// out a branch from the remote. Config set using this approach will override
// any config defined within existing git config files. Config must be
// provided as key value pairs, mismatched config will result in an
// [ErrMissingConfigValue] error. Any invalid paths will result in an
// [ErrInvalidConfigPath] error
func WithCheckoutConfig(kv ...string) CheckoutOption {
	return func(opts *checkoutOptions) {
		opts.Config = trim(kv...)
	}
}

// Checkout will attempt to checkout a branch with the given name. If the branch
// does not exist, it is created at the current working tree reference (or commit),
// and then switched to. If the branch does exist, then switching to it restores
// all working tree files
func (c *Client) Checkout(branch string, opts ...CheckoutOption) (string, error) {
	options := &checkoutOptions{}
	for _, opt := range opts {
		opt(options)
	}

	cfg, err := ToInlineConfig(options.Config...)
	if err != nil {
		return "", err
	}

	// Query the repository for all existing branches, both local and remote.
	// If a pull hasn't been done, there is a chance that an expected
	// remote branch will not be tracked
	out, err := c.exec("git branch --all --format='%(refname:short)'")
	if err != nil {
		return out, err
	}

	var buf strings.Builder
	buf.WriteString("git")

	if len(cfg) > 0 {
		buf.WriteString(" ")
		buf.WriteString(strings.Join(cfg, " "))
	}
	buf.WriteString(" checkout ")

	for _, ref := range strings.Split(out, "\n") {
		if strings.HasSuffix(ref, branch) {
			return c.exec(buf.String() + branch)
		}
	}

	buf.WriteString(" -b ")
	return c.exec(buf.String() + branch)
}
