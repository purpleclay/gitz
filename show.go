package git

import (
	"strings"
	"time"
)

const (
	dateFormat   = "Mon Jan _2 15:04:05 2006 -0700"
	commitIndent = "    "
)

// BlobDetails contains details about a specific blob within a repository.
type BlobDetails struct {
	// Diff contains the blobs contents
	Diff string

	// Ref contains the unique identifier associated with the blob
	Ref string
}

// CommitDetails contains details about a specific commit within a repository.
type CommitDetails struct {
	// Author represents a person who originally created the files
	// within the repository
	Author Person

	// AuthorDate contains the date and time of when the author
	// originally created the files within the repository
	AuthorDate time.Time

	// Committer represents a person who changed any existing files
	// within the repository
	Committer Person

	// CommitterDate contains the date and time of when the committer
	// changed any existing files within the repository
	CommitterDate time.Time

	// Message contains the message associated with the commit
	Message string

	// Ref contains the unique identifier associated with the commit
	Ref string

	// Signature contains details of the verified GPG signature
	Signature *Signature
}

// TagAnnotation contains details about an annotation associated with a tag
// within a repository.
type TagAnnotation struct {
	// Tagger represents a person who created the tag
	Tagger Person

	// TaggerDate contains the date and time of when the tagger created
	// the tag within the repository
	TaggerDate time.Time

	// Message contains the annotated message associated with the tag
	Message string
}

// TagDetails contains details about a specific tag within a repository.
type TagDetails struct {
	// Annotation contains optional details about the annotation associated
	// with the tag
	Annotation *TagAnnotation

	// Commit contains details about the associated commit
	Commit CommitDetails

	// Ref contains the unique identifier associated with the tag
	Ref string
}

// TreeDetails contains details about a specific tree within a repository.
type TreeDetails struct {
	// Entries contains the file and directory listing within a tree
	Entries []string

	// Ref contains the unique identifier associated with the tree
	Ref string
}

// Person represents a human that has performed an interaction against
// a repository.
type Person struct {
	// Name of the person
	Name string

	// Email address associated with the person
	Email string
}

// ShowBlobs retrieves details about any number of blobs within the current
// repository (working directory).
func (c *Client) ShowBlobs(refs ...string) (map[string]BlobDetails, error) {
	details := map[string]BlobDetails{}
	for _, ref := range refs {
		diff, err := c.Exec("git show --no-color " + ref)
		if err != nil {
			return nil, err
		}
		details[ref] = BlobDetails{Diff: diff, Ref: ref}
	}

	return details, nil
}

// ShowCommits retrieves details about any number of commits within the current
// repository (working directory).
func (c *Client) ShowCommits(refs ...string) (map[string]CommitDetails, error) {
	details := map[string]CommitDetails{}
	for _, ref := range refs {
		out, err := c.Exec("git show --no-color -s --show-signature --format=fuller " + ref)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(out, "commit") {
			commit := parseCommit(out)
			commit.Ref = ref

			details[ref] = commit
		}
	}

	return details, nil
}

func parseCommit(str string) CommitDetails {
	str, _ = line()(str)
	var signature *Signature
	if strings.HasPrefix(str, "gpg:") {
		var gpg string
		str, gpg = until("Author:")(str)
		signature = parseSignature(gpg)
	}
	str, pair := separatedPair(tag("Author:"), ws(), until("AuthorDate:"))(str)
	author := parsePerson(pair[1])

	str, pair = separatedPair(tag("AuthorDate:"), ws(), until("Commit:"))(str)
	authorDate, _ := time.Parse(dateFormat, chompCRLF(pair[1]))

	str, pair = separatedPair(tag("Commit:"), ws(), until("CommitDate:"))(str)
	committer := parsePerson(pair[1])

	str, pair = separatedPair(tag("CommitDate:"), ws(), until("\n"))(str)
	committerDate, _ := time.Parse(dateFormat, chompCRLF(pair[1]))

	return CommitDetails{
		Author:        author,
		AuthorDate:    authorDate,
		Committer:     committer,
		CommitterDate: committerDate,
		Signature:     signature,
		Message:       strings.TrimSpace(chompIndent(commitIndent, str)),
	}
}

// ShowTags retrieves details about any number of tags within the current
// repository (working directory).
func (c *Client) ShowTags(refs ...string) (map[string]TagDetails, error) {
	details := map[string]TagDetails{}
	for _, ref := range refs {
		show, err := c.Exec("git show --no-color -s --show-signature --format=fuller " + ref)
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(show, "tag") {
			str, _ := until("Tagger:")(show)

			str, pair := separatedPair(tag("Tagger:"), ws(), until("TaggerDate:"))(str)
			tagger := parsePerson(pair[1])

			str, pair = separatedPair(tag("TaggerDate:"), ws(), until("\n"))(str)
			taggerDate, _ := time.Parse(dateFormat, chompCRLF(pair[1]))

			str, _ = takeUntil(alphanumeric)(str)
			str, message := until("commit")(str)
			if i := strings.Index(message, "-----BEGIN PGP SIGNATURE-----"); i != -1 {
				message = message[:i]
			}
			message = strings.TrimSpace(message)

			details[ref] = TagDetails{
				Annotation: &TagAnnotation{
					Tagger:     tagger,
					TaggerDate: taggerDate,
					Message:    message,
				},
				Commit: parseCommit(str),
				Ref:    ref,
			}
		} else if strings.HasPrefix(show, "commit") {
			details[ref] = TagDetails{
				Commit: parseCommit(show),
				Ref:    ref,
			}
		}
	}

	return details, nil
}

// ShowTrees retrieves details about any number of trees within the current
// repository (working directory).
func (c *Client) ShowTrees(refs ...string) (map[string]TreeDetails, error) {
	details := map[string]TreeDetails{}
	for _, ref := range refs {
		tree, err := c.Exec("git show --no-color " + ref)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(tree, "tree") {
			// Trim the first two lines that represent the tree header
			tree = tree[strings.Index(tree, "\n\n")+2:]
			details[ref] = TreeDetails{
				Entries: strings.Split(tree, "\n"),
				Ref:     ref,
			}
		}
	}

	return details, nil
}
