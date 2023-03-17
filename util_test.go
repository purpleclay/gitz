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

package git_test

import (
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	tests := []struct {
		name     string
		strs     []string
		expected []string
	}{
		{
			name:     "RemovesWhitespace",
			strs:     []string{" how", "are   ", "  you  "},
			expected: []string{"how", "are", "you"},
		},
		{
			name:     "RemovesEmptyStrings",
			strs:     []string{"hello", " ", "    ", "there"},
			expected: []string{"hello", "there"},
		},
		{
			name:     "ReturnsEmptySlice",
			strs:     []string{" ", "    "},
			expected: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trimmed := git.Trim(tt.strs...)
			assert.ElementsMatch(t, tt.expected, trimmed)
		})
	}
}

func TestTrimAndPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		strs     []string
		expected []string
	}{
		{
			name:     "RemovesWhitespaceAppendsPrefix",
			prefix:   "refs/tags/",
			strs:     []string{"  0.1.0", "  0.2.0  ", "0.3.0  "},
			expected: []string{"refs/tags/0.1.0", "refs/tags/0.2.0", "refs/tags/0.3.0"},
		},
		{
			name:     "RemovesEmptyStrings",
			prefix:   "job-",
			strs:     []string{"#1", " ", "", "#2"},
			expected: []string{"job-#1", "job-#2"},
		},
		{
			name:     "ReturnsEmptySlice",
			prefix:   "--sort=",
			strs:     []string{" ", "    "},
			expected: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trimmed := git.TrimAndPrefix(tt.prefix, tt.strs...)
			assert.ElementsMatch(t, tt.expected, trimmed)
		})
	}
}
