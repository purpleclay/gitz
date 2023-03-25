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

package scan

import "bytes"

// PrefixedLines is a split function for a [bufio.Scanner] that returns
// each block of text, stripped of both the prefix marker and any leading
// and trailing whitespace. If no prefix is detected, the original text
// will be treated as a single block of text, with any leading and trailing
// whitespace stripped
func PrefixedLines(prefix byte) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.Index(data, []byte{'\n', prefix}); i >= 0 {
			return i + 1, eat(prefix, data[:i]), nil
		}

		if atEOF {
			return len(data), eat(prefix, data), nil
		}

		return 0, nil, nil
	}
}

func eat(prefix byte, data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	i := 0
	if data[i] == prefix {
		i++
	}

	return bytes.TrimSpace(data[i:])
}
