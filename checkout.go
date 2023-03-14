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

// Checkout ...
func (c *Client) Checkout(branch string) (string, error) {
	// Query the repository for all existing branches, both local and remote.
	// If a pull hasn't been done, there is a chance that an expected
	// remote branch will not be tracked
	out, err := exec("git branch --all --format='%(refname:short)'")
	if err != nil {
		return out, err
	}

	for _, ref := range strings.Split(out, "\n") {
		if strings.HasSuffix(ref, branch) {
			return exec(fmt.Sprintf("git checkout %s", branch))
		}
	}

	return exec(fmt.Sprintf("git checkout -b %s", branch))
}
