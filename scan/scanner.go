package scan

import (
	"bytes"
)

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
	i := 0
	if i < len(data) && data[i] == prefix {
		i++
	}

	return bytes.TrimSpace(data[i:])
}

// DiffLines is a split function for a [bufio.Scanner] that splits a git diff output
// into multiple blocks of text, each prefixed by the diff --git marker. Each block
// of text will be stripped of any leading and trailing whitespace. If the git diff
// marker isn't detected, the entire block of text is returned, with any leading and
// trailing whitespace stripped
func DiffLines() func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	prefix := []byte("\ndiff --git")

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.Index(data, prefix); i >= 0 {
			return i + 1, bytes.TrimSpace(data[:i]), nil
		}

		if atEOF {
			return len(data), bytes.TrimSpace(data), nil
		}

		return 0, nil, nil
	}
}
