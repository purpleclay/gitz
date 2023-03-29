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
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Config ...
func (c *Client) Config(path string) (string, error) {
	var cmd strings.Builder
	cmd.WriteString("git config --get ")
	cmd.WriteString(path)

	return exec(cmd.String())
}

// ConfigL ...
func (c *Client) ConfigL(paths ...string) (map[string]string, error) {
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

// ConfigSet ...
func (c *Client) ConfigSet(path, value string) error {
	var cmd strings.Builder
	cmd.WriteString("git config --add ")
	cmd.WriteString(fmt.Sprintf("%s '%s'", path, value))

	_, err := exec(cmd.String())
	return err
}

// ConfigSetL ...
func (c *Client) ConfigSetL(pairs ...string) error {
	if len(pairs) == 0 {
		return nil
	}

	if len(pairs)%2 != 0 {
		return errors.New("uneven pairs provided, not enough values to paths") // print out the offending path
	}

	for i := 0; i < len(pairs); i += 2 {
		fmt.Println(pairs[i])
		if !ValidConfigPath(pairs[i]) {
			return errors.New("path is not valid") // must be made up of ...
		}
	}

	for i := 0; i < len(pairs); i += 2 {
		if err := c.ConfigSet(pairs[i], pairs[i+1]); err != nil {
			return err
		}
	}

	return nil
}

// ValidConfigPath ...
func ValidConfigPath(path string) bool {
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 || lastDot == len(path)-1 {
		return false
	}

	for i, c := range path {
		if i == lastDot+1 && !unicode.IsLetter(c) {
			return false
		}

		if unicode.IsDigit(c) || unicode.IsLetter(c) || c == '.' {
			continue
		}

		return false
	}

	return true
}
