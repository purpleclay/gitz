package git

import (
	"fmt"
	"strings"
)

// CommitOption provides a way for setting specific options during a commit
// operation. Each supported option can customize the way the commit is
// created against the current repository (working directory)
type CommitOption func(*commitOptions)

type commitOptions struct {
	AllowEmpty    bool
	Config        []string
	ForceNoSigned bool
	Signed        bool
	SigningKey    string
}

// WithAllowEmpty allows a commit to be created without having to track
// any changes. This bypasses the default protection by git, preventing
// a commit from having the exact same tree as its parent
func WithAllowEmpty() CommitOption {
	return func(opts *commitOptions) {
		opts.AllowEmpty = true
	}
}

// WithCommitConfig allows temporary git config to be set during the
// execution of the commit. Config set using this approach will override
// any config defined within existing git config files. Config must be
// provided as key value pairs, mismatched config will result in an
// [ErrMissingConfigValue] error. Any invalid paths will result in an
// [ErrInvalidConfigPath] error
func WithCommitConfig(kv ...string) CommitOption {
	return func(opts *commitOptions) {
		opts.Config = trim(kv...)
	}
}

// WithGpgSign will create a GPG-signed commit using the GPG key associated
// with the committers email address. Overriding this behavior is possible
// through the user.signingkey config setting. This option does not need
// to be explicitly called if the commit.gpgSign config setting is set to
// true
func WithGpgSign() CommitOption {
	return func(opts *commitOptions) {
		opts.Signed = true
	}
}

// WithGpgSigningKey will create a GPG-signed commit using the provided GPG
// key ID, overridding any default GPG key set by the user.signingKey git
// config setting
func WithGpgSigningKey(key string) CommitOption {
	return func(opts *commitOptions) {
		opts.Signed = true
		opts.SigningKey = strings.TrimSpace(key)
	}
}

// WithNoGpgSign ensures the created commit will not be GPG signed
// regardless of the value assigned to the repositories commit.gpgSign
// git config setting
func WithNoGpgSign() CommitOption {
	return func(opts *commitOptions) {
		opts.ForceNoSigned = true
	}
}

// Commit a snapshot of changes within the current repository (working directory)
// and describe those changes with a given log message. Commit behavior can be
// customized through the use of options
func (c *Client) Commit(msg string, opts ...CommitOption) (string, error) {
	options := &commitOptions{}
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
	buf.WriteString(" commit")

	if options.AllowEmpty {
		buf.WriteString(" --allow-empty")
	}

	if options.Signed {
		buf.WriteString(" -S")
	}

	if options.SigningKey != "" {
		buf.WriteString(" --gpg-sign=" + options.SigningKey)
	}

	if options.ForceNoSigned {
		buf.WriteString(" --no-gpg-sign")
	}

	buf.WriteString(fmt.Sprintf(" -m '%s'", msg))
	return c.Exec(buf.String())
}

// CommitVerification contains details about a GPG signed commit
type CommitVerification struct {
	// Author represents a person who originally created the files
	// within the repository
	Author Person

	// Committer represents a person who changed any existing files
	// within the repository
	Committer Person

	// Hash contains the unique identifier associated with the commit
	Hash string

	// Message contains the message associated with the commit
	Message string

	// Signature contains details of the verified GPG signature
	Signature *Signature
}

// VerifyCommit validates that a given commit has a valid GPG signature
// and returns details about that signature
func (c *Client) VerifyCommit(hash string) (*CommitVerification, error) {
	out, err := c.Exec("git verify-commit -v " + hash)
	if err != nil {
		return nil, err
	}

	out, _ = until("author ")(out)
	out, pair := separatedPair(tag("author "), ws(), until("committer "))(out)
	author := parsePerson(pair[1])

	out, pair = separatedPair(tag("committer "), ws(), takeUntil(lineEnding))(out)
	committer := parsePerson(pair[1])
	out, _ = line()(out)

	out, mesage := until("gpg: ")(out)

	return &CommitVerification{
		Author:    author,
		Committer: committer,
		Hash:      hash,
		Message:   strings.TrimSpace(mesage),
		Signature: parseSignature(out),
	}, nil
}
