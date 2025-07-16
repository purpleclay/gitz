package git

import (
	"fmt"
	"strings"
	"unicode"
)

// InvalidConfigPathError is raised when a config setting is to be accessed
// with an invalid config path.
type InvalidConfigPathError struct {
	// Path to the config setting.
	Path string

	// Position of the first offending character within the path.
	Position int

	// Reason why the path is invalid.
	Reason string
}

// Error returns a friendly formatted message of the current error.
func (e InvalidConfigPathError) Error() string {
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

// MissingConfigValueError is raised when a git config path does not
// have a corresponding value.
type MissingConfigValueError struct {
	// Path to the config setting
	Path string
}

// Error returns a friendly formatted message of the current error.
func (e MissingConfigValueError) Error() string {
	return fmt.Sprintf("config paths mismatch. path: %s is missing a corresponding value", e.Path)
}

// Config attempts to retrieve all git config for the current repository.
// A map is returned containing each config item and its corresponding
// latest value. Values are resolved from local, system and global config.
func (c *Client) Config() (map[string]string, error) {
	cfg, err := c.Exec("git config --list")
	if err != nil {
		return nil, err
	}

	values := map[string]string{}

	lines := strings.Split(cfg, "\n")
	for _, line := range lines {
		pos := strings.Index(line, "=")
		values[line[:pos]] = line[pos+1:]
	}

	return values, nil
}

// ConfigL attempts to query a batch of local git config settings for
// their values. If multiple values have been set for any config item,
// all are returned, ordered by most recent value first. A partial batch
// is never returned, all config settings must exist.
func (c *Client) ConfigL(paths ...string) (map[string][]string, error) {
	return c.configQuery("local", paths...)
}

func (c *Client) configQuery(location string, paths ...string) (map[string][]string, error) {
	values := map[string][]string{}

	var cmd strings.Builder
	for _, path := range paths {
		cmd.WriteString("git config ")
		cmd.WriteString("--" + location)
		cmd.WriteString(" --get-all ")
		cmd.WriteString(path)

		cfg, err := c.Exec(cmd.String())
		if err != nil {
			return nil, err
		}
		cmd.Reset()

		v := reverse(strings.Split(cfg, "\n")...)
		values[path] = v
	}

	return values, nil
}

// ConfigG attempts to query a batch of global git config settings for
// their values. If multiple values have been set for any config item,
// all are returned, ordered by most recent value first. A partial batch
// is never returned, all config settings must exist.
func (c *Client) ConfigG(paths ...string) (map[string][]string, error) {
	return c.configQuery("global", paths...)
}

// ConfigS attempts to query a batch of system git config settings for
// their values. If multiple values have been set for any config item,
// all are returned, ordered by most recent value first. A partial batch
// is never returned, all config settings must exist.
func (c *Client) ConfigS(paths ...string) (map[string][]string, error) {
	return c.configQuery("system", paths...)
}

// ConfigSetL attempts to batch assign values to a group of local git
// config settings. If any setting exists, a new line is added to the
// local git config, effectively assigning multiple values to the same
// setting. Basic validation is performed to minimize the possibility
// of a partial batch update.
func (c *Client) ConfigSetL(pairs ...string) error {
	return c.configSet("local", pairs...)
}

func (c *Client) configSet(location string, pairs ...string) error {
	if len(pairs) == 0 {
		return nil
	}

	if err := checkConfig(pairs); err != nil {
		return err
	}

	var cmd strings.Builder
	for i := 0; i < len(pairs); i += 2 {
		cmd.WriteString("git config ")
		cmd.WriteString("--" + location)
		cmd.WriteString(" --add ")
		cmd.WriteString(fmt.Sprintf("%s '%s'", pairs[i], pairs[i+1]))

		if _, err := c.Exec(cmd.String()); err != nil {
			return err
		}
		cmd.Reset()
	}

	return nil
}

// ConfigSetG attempts to batch assign values to a group of global git
// config settings. If any setting exists, a new line is added to the
// local git config, effectively assigning multiple values to the same
// setting. Basic validation is performed to minimize the possibility
// of a partial batch update.
func (c *Client) ConfigSetG(pairs ...string) error {
	return c.configSet("global", pairs...)
}

// ConfigSetS attempts to batch assign values to a group of system git
// config settings. If any setting exists, a new line is added to the
// local git config, effectively assigning multiple values to the same
// setting. Basic validation is performed to minimize the possibility
// of a partial batch update.
func (c *Client) ConfigSetS(pairs ...string) error {
	return c.configSet("system", pairs...)
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
		return InvalidConfigPathError{
			Path:     path,
			Position: lastDot,
			Reason:   "dot separator is missing or is the last character",
		}
	}

	for i, c := range path {
		if i == lastDot+1 && !unicode.IsLetter(c) {
			return InvalidConfigPathError{
				Path:     path,
				Position: i,
				Reason:   "first character after final dot must be a letter [a-zA-Z]",
			}
		}

		if unicode.IsDigit(c) || unicode.IsLetter(c) || c == '.' {
			continue
		}

		return InvalidConfigPathError{
			Path:     path,
			Position: i,
			Reason:   "non alphanumeric character detected [a-zA-Z0-9]",
		}
	}

	return nil
}

func checkConfig(pairs []string) error {
	if len(pairs)%2 != 0 {
		return MissingConfigValueError{Path: pairs[len(pairs)-1]}
	}

	for i := 0; i < len(pairs); i += 2 {
		if err := CheckConfigPath(pairs[i]); err != nil {
			return err
		}
	}

	return nil
}

// ToInlineConfig converts a series of config settings from path value notation
// into the corresponding inline config notation compatible with git commands
//
//	"user.name", "penguin" => []string{"-c user.name='penguin'"}
func ToInlineConfig(pairs ...string) ([]string, error) {
	if len(pairs) == 0 {
		return nil, nil
	}

	if err := checkConfig(pairs); err != nil {
		return nil, err
	}

	cfg := make([]string, 0, len(pairs)%2)
	for i := 0; i < len(pairs); i += 2 {
		cfg = append(cfg, fmt.Sprintf("-c %s='%s'", pairs[i], pairs[i+1]))
	}

	return cfg, nil
}
