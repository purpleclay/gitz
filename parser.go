package git

import (
	"strings"
	"unicode"
)

type parser func(string) (string, []string)

type combinator func(string) (string, string)

func separatedPair(first, sep, second combinator) parser {
	return func(s string) (string, []string) {
		out := make([]string, 0, 2)
		str, ret := first(s)
		out = append(out, ret)
		str, _ = sep(str)
		str, ret = second(str)
		out = append(out, ret)

		return str, out
	}
}

func tag(tag string) combinator {
	return func(s string) (string, string) {
		if strings.HasPrefix(s, tag) {
			return s[len(tag):], tag
		}
		return s, ""
	}
}

func ws() combinator {
	return func(s string) (string, string) {
		for i, c := range s {
			if !unicode.IsSpace(c) {
				return s[i:], s[:i]
			}
		}
		return s, ""
	}
}

func until(delim string) combinator {
	return func(s string) (string, string) {
		if i := strings.Index(s, delim); i > -1 {
			return s[i:], s[:i]
		}
		return s, ""
	}
}

func line() combinator {
	return func(s string) (string, string) {
		if i := strings.Index(s, "\n"); i > 0 {
			j := i
			if j > 1 && s[j-1] == '\r' {
				j--
			}

			if len(s) == i {
				return "", s[:j]
			}
			return s[i+1:], s[:j]
		}
		return s, ""
	}
}

type condition func(string) int

func alphanumeric(str string) int {
	for i, b := range str {
		if unicode.IsLetter(b) || unicode.IsNumber(b) {
			return i
		}
	}
	return -1
}

func lineEnding(str string) int {
	for i, b := range str {
		if b == '\r' || b == '\n' {
			return i
		}
	}
	return -1
}

func takeUntil(cond condition) combinator {
	return func(s string) (string, string) {
		if i := cond(s); i != -1 {
			return s[i:], s[:i]
		}
		return s, s
	}
}
