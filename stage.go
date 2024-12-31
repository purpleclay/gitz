package git

import "strings"

// StageOption provides a way for setting specific options during a stage
// operation. Each supported option can customize the way files are staged
// within the current repository (working directory)
type StageOption func(*stageOptions)

type stageOptions struct {
	PathSpecs []string
}

// WithPathSpecs permits a series of [PathSpecs] (or globs) to be defined
// that will stage any matching files within the current repository
// (working directory). Paths to files and folders are relative to the
// root of the repository. All leading and trailing whitespace will be
// trimmed from the file paths, allowing empty paths to be ignored
//
// [PathSpecs]: https://git-scm.com/docs/gitglossary#Documentation/gitglossary.txt-aiddefpathspecapathspec
func WithPathSpecs(specs ...string) StageOption {
	return func(opts *stageOptions) {
		opts.PathSpecs = trim(specs...)
	}
}

// Stage changes to any file or folder within the current repository
// (working directory) ready for inclusion in the next commit. Options
// can be provided to further configure stage semantics. By default,
// all changes will be staged ready for the next commit.
func (c *Client) Stage(opts ...StageOption) (string, error) {
	options := &stageOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Build command based on the provided options
	var stageCmd strings.Builder
	stageCmd.WriteString("git add ")

	if len(options.PathSpecs) > 0 {
		stageCmd.WriteString("--")
		for _, spec := range options.PathSpecs {
			stageCmd.WriteString(" ")
			stageCmd.WriteString(spec)
		}
	} else {
		stageCmd.WriteString("--all")
	}

	return c.Exec(stageCmd.String())
}

// Staged retrieves a list of all currently staged file changes within the
// current repository
func (c *Client) Staged() ([]string, error) {
	diff, err := c.Exec("git diff --staged --name-only")
	if err != nil {
		return nil, err
	}

	if diff == "" {
		return nil, nil
	}

	return strings.Split(diff, "\n"), nil
}
