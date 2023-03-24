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

// PrefixedLines ...
func PrefixedLines(prefix byte) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Check for the existence of the expected marker
		if i := bytes.IndexByte(data, prefix); i != 0 {
			return 0, nil, nil
		}

		if i := bytes.Index(data, []byte{'\n', prefix}); i >= 0 {
			return i + 1, dropCR(eatWS(data[1:i])), nil
		}

		if atEOF {
			return len(data), dropCR(eatWS(data[1:])), nil
		}

		return 0, nil, nil
	}
}

func eatWS(data []byte) []byte {
	i := 0
	for ; i < len(data); i++ {
		if data[i] != ' ' {
			break
		}
	}

	return data[i:]
}

// shamelessly copied from the bufio package
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}

	return data
}
