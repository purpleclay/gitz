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

// CommitOption provides a way for setting specific options during a commit
// operation. Each supported option can customize the way the commit is
// created against the current repository (working directory)
type CommitOption func(*commitOptions)

type commitOptions struct {
	AllowEmpty    bool
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

	var commitCmd strings.Builder
	commitCmd.WriteString("git commit")

	if options.AllowEmpty {
		commitCmd.WriteString(" --allow-empty")
	}

	if options.Signed {
		commitCmd.WriteString(" -S")
	}

	if options.SigningKey != "" {
		commitCmd.WriteString(" --gpg-sign=" + options.SigningKey)
	}

	if options.ForceNoSigned {
		commitCmd.WriteString(" --no-gpg-sign")
	}

	commitCmd.WriteString(fmt.Sprintf(" -m '%s'", msg))
	return c.exec(commitCmd.String())
}

const (
	authorPrefix    = "author "
	committerPrefix = "committer "
	emailEnd        = '>'
)

// CommitVerification contains details about a GPG signed commit
type CommitVerification struct {
	Sha         string
	Author      Author
	Committer   Author
	Fingerprint string
	SignedBy    *Author
}

// VerifyCommit validates that a given commit has a valid GPG signature
// and returns details about that signature
func (c *Client) VerifyCommit(sha string) (*CommitVerification, error) {
	out, err := c.exec("git verify-commit -v " + sha)
	if err != nil {
		return nil, err
	}

	author := chompUntil(out[strings.Index(out, authorPrefix)+len(authorPrefix):], emailEnd)
	committer := chompUntil(out[strings.Index(out, committerPrefix)+len(committerPrefix):], emailEnd)
	fingerprint := chompCRLF(out[strings.Index(out, fingerprintPrefix)+len(fingerprintPrefix):])

	var signedByAuthor *Author
	if strings.Contains(out, signedByPrefix) {
		signedBy := chompUntil(out[strings.Index(out, signedByPrefix)+len(signedByPrefix):], '"')
		author := parseAuthor(signedBy)
		signedByAuthor = &author
	}

	return &CommitVerification{
		Sha:         sha,
		Author:      parseAuthor(author),
		Committer:   parseAuthor(committer),
		Fingerprint: fingerprint,
		SignedBy:    signedByAuthor,
	}, nil
}
