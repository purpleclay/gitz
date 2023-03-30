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
	"unicode"
)

// ErrInvalidConfigPath is raised when a config setting is to be accessed
// with an invalid config path
type ErrInvalidConfigPath struct {
	// Path to the config setting
	Path string

	// Position of the first offending character within the path
	Position int

	// Reason why the path is invalid
	Reason string
}

// Error returns a friendly formatted message of the current error
func (e ErrInvalidConfigPath) Error() string {
	var buf strings.Builder
	if e.Position == -1 {
		buf.WriteString(e.Path)
	} else {
		buf.WriteString(e.Path[:e.Position])
		buf.WriteString(fmt.Sprintf("|%c|", e.Path[e.Position]))
		if e.Position != len(e.Path)-1 {
			buf.WriteString(e.Path[e.Position+1:])
		}
	}

	return fmt.Sprintf("path: %s invalid as %s", buf.String(), e.Reason)
}

// Config attempts to query a local git config setting for its value.
// If multiple values have been set, all are returned, ordered by the
// most recent value first
func (c *Client) Config(path string) ([]string, error) {
	var cmd strings.Builder
	cmd.WriteString("git config --get-all ")
	cmd.WriteString(path)

	// TODO: switch to parsing the result into separate values and reversing the slice

	return exec(cmd.String())
}

// ConfigL attempts to query a batch of local git config settings for
// their values. If multiple values have been set for any config item,
// all are returned, ordered by most recent value first. A partial batch
// is never returned, all config settings must exist
func (c *Client) ConfigL(paths ...string) (map[string][]string, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	cfg := map[string]string{}
	for _, path := range paths {
		v, err := c.Config(path)
		if err != nil {
			return nil, err
		}

		cfg[path] = v
	}

	return cfg, nil
}

// ConfigSet attempts to assign a value to a local git config setting.
// If the setting already exists, a new line is added to the local git
// config, effectively assigning multiple values to the same setting
func (c *Client) ConfigSet(path, value string) error {
	var cmd strings.Builder
	cmd.WriteString("git config --add ")
	cmd.WriteString(fmt.Sprintf("%s '%s'", path, value))

	_, err := exec(cmd.String())
	return err
}

// ConfigSetL attempts to batch assign values to a group of local git
// config settings. If any setting exists, a new line is added to the
// local git config, effectively assigning multiple values to the same
// setting. Basic validation is performed to minimize the possibility
// of a partial batch update
func (c *Client) ConfigSetL(pairs ...string) error {
	if len(pairs) == 0 {
		return nil
	}

	if len(pairs)%2 != 0 {
		return fmt.Errorf("config paths mismatch. path: %s is missing a corresponding value", pairs[len(pairs)-1])
	}

	for i := 0; i < len(pairs); i += 2 {
		if err := CheckConfigPath(pairs[i]); err != nil {
			return err
		}
	}

	for i := 0; i < len(pairs); i += 2 {
		if err := c.ConfigSet(pairs[i], pairs[i+1]); err != nil {
			return err
		}
	}

	return nil
}

// CheckConfigPath performs rudimentary checks to ensure the config path
// conforms to the git config specification. A config path is invalid if:
//
//   - No dot separator exists, or the last character is a dot separator
//   - First character after the last dot separator is not a letter
//   - Path contains non-alphanumeric characters
func CheckConfigPath(path string) error {
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 || lastDot == len(path)-1 {
		return ErrInvalidConfigPath{
			Path:     path,
			Position: lastDot,
			Reason:   "dot separator is missing or is the last character",
		}
	}

	for i, c := range path {
		if i == lastDot+1 && !unicode.IsLetter(c) {
			return ErrInvalidConfigPath{
				Path:     path,
				Position: i,
				Reason:   "first character after final dot must be a letter [a-zA-Z]",
			}
		}

		if unicode.IsDigit(c) || unicode.IsLetter(c) || c == '.' {
			continue
		}

		return ErrInvalidConfigPath{
			Path:     path,
			Position: i,
			Reason:   "non alphanumeric character detected [a-zA-Z0-9]",
		}
	}

	return nil
}
