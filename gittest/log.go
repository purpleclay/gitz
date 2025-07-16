package gittest

import (
	"bufio"
	"strings"

	"github.com/purpleclay/gitz/scan"
)

// LogEntry represents a single log entry from the history
// of a git repository.
type LogEntry struct {
	// Hash contains the unique identifier associated with the commit.
	Hash string

	// AbbrevHash contains the seven character abbreviated commit hash.
	AbbrevHash string

	// Message contains the log message associated with the commit.
	Message string

	// Tags contains all tag references that are associated
	// with the current commit.
	Tags []string

	// Branches contains the name of all branches (local and remote)
	// that are associated with the current commit.
	Branches []string

	// IsTrunk identifies if the current log entry has a reference
	// to the default branch.
	IsTrunk bool

	// HeadPointerRef contains the name of the branch the HEAD of the
	// repository points to.
	HeadPointerRef string
}

// ParseLog will attempt to parse a log extract from a given repository
// into a series of commits, branches and tags. The log will be returned
// in the chronological order provided. The parser is designed to not
// error and parses each line with best endeavors. Multiple log formats
// are supported:
//
// 1. A condensed single line log format:
//
//	(HEAD -> new-feature, origin/new-feature) pass tests
//	write tests for new feature
//	(tag: 0.2.0, tag: v1, main, origin/main) feat: improve existing cli documentation
//	docs: create initial mkdocs material documentation
//	(tag: 0.1.0) feat: add secondary cli command to support filtering of results
//	feat: scaffold initial cli and add first command
//
// This is the equivalent to the format produced using the git command:
//
//	git log --pretty='format:%d %s'
//
// 2. A multi-line commit format, which is denotoed with the presence of the prefix
// marker >. The parser will switch between parsing modes if it detects the existence
// of this marker on the first character of the provided log extract:
//
//	> (HEAD -> new-feature, origin/new-feature) pass tests
//	> write tests for new feature
//	> (tag: 0.2.0, tag: v1, main, origin/main) feat: improve existing cli documentation
//
// This is the equivalent to the format produced using the git command:
//
//	git log --pretty='format:> %d %s%+b%-N'
//
// 3. A log containing an optional leading forty character hash. Can be used
// in conjunction with both single line and multi-line formats:
//
//	> b0d5429b967b9af0a0805fc2981b4420e10be38d (HEAD -> new-feature, origin/new-feature) pass tests
//	> 58d708cb071df97e2561903aadcd4129419e9631 write tests for new feature
//	> 4edd1a7e492aeeaf2a97ad57433e236bc72e1d93 (tag: 0.2.0, tag: v1, main, origin/main) feat: improve existing cli documentation
//
// This is the equivalent to the format produced using the git command:
//
//	git log --pretty='format:> %H %d %s%+b%-N'
//
// [%m]: https://git-scm.com/docs/git-log#Documentation/git-log.txt-emmem
func ParseLog(log string) []LogEntry {
	if log == "" {
		return nil
	}

	entries := make([]LogEntry, 0)
	scanner := bufio.NewScanner(strings.NewReader(log))

	// Detect if the log requires multi-line parsing by checking for the git marker > (%m)
	if log[0] == '>' {
		scanner.Split(scan.PrefixedLines('>'))
	} else {
		scanner.Split(bufio.ScanLines)
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		var hash string
		var abbrevHash string
		if hash, line = chompHash(line); len(hash) > 0 {
			abbrevHash = hash[:7]
			// Ensure any leading whitespace is removed before proceeding
			line = strings.TrimSpace(line)
		}

		entry := LogEntry{
			Hash:       hash,
			AbbrevHash: abbrevHash,
			Message:    line,
		}
		if strings.HasPrefix(line, "(") {
			// Cut based on the first occurrence of a closing parentheses, if one doesn't
			// exist, then append the line as a raw log entry
			refNames, msg, found := strings.Cut(line, ") ")
			if !found {
				goto append
			}
			entry.Message = msg

			// Process the comma separated list of ref names, preceding the commit message
			for _, ref := range strings.Split(refNames[1:], ",") {
				cleanedRef := strings.TrimSpace(ref)

				if cleanedRef == "" {
					continue
				}

				if strings.HasPrefix(cleanedRef, "tag: ") {
					entry.Tags = append(entry.Tags, strings.TrimPrefix(cleanedRef, "tag: "))
				} else {
					entry.Branches = append(entry.Branches, cleanedRef)

					// Detect the existence of the default branch
					if cleanedRef == DefaultBranch {
						entry.IsTrunk = true
					}
				}
			}
		}

	append:
		entries = append(entries, entry)
	}

	// Determine if the first log entry contains a HEAD pointer reference
	for _, branch := range entries[0].Branches {
		if hasBranchPrefix(branch, "HEAD->", "HEAD ->") {
			if _, pointer, found := strings.Cut(branch, "->"); found {
				head := strings.TrimSuffix(pointer, DefaultBranch)
				entries[0].HeadPointerRef = strings.TrimSpace(head)
				break
			}
		}
	}

	return entries
}

func hasBranchPrefix(branch string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(branch, prefix) {
			return true
		}
	}
	return false
}

func chompHash(str string) (string, string) {
	if len(str) < 40 {
		return "", str
	}

	hash := str[:40]
	for _, b := range []byte(hash) {
		if (b < '0' || b > '9') && (b < 'a' || b > 'f') && (b < 'A' || b > 'F') {
			return "", str
		}
	}
	return hash, str[40:]
}
