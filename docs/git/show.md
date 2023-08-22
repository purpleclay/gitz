---
icon: material/text-search
title: Inspect an object within a repository
description: Retrieve details about a specific object from within a repository
---

# Inspect an object within a repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-show)

Retrieve detailed information about an object within a repository by its unique reference.

## Inspect a tag

Detailed information about a tag, and its associated commit, can be retrieved from a repository by passing its reference to `ShowTags`. The GPG signature of the commit is also retrieved if present.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "time"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Querying a tag from the gpg-import project

    tags, _ := client.ShowTags("0.3.2")

    tag := tags[0]
    if tag.Annotation != nil {
        fmt.Printf("Tagger:      %s <%s>\n",
            tag.Annotation.Tagger.Name, tag.Annotation.Tagger.Email)
        fmt.Printf("TaggerDate:  %s\n",
            tag.Annotation.TaggerDate.Format(time.RubyDate))
        fmt.Printf("Message:     %s\n\n", tag.Annotation.Message)
    }

    fmt.Printf("Author:      %s <%s>\n",
        tag.Commit.Author.Name, tag.Commit.Author.Email)
    fmt.Printf("AuthorDate:  %s\n",
        tag.Commit.AuthorDate.Format(time.RubyDate))
    fmt.Printf("Commit:      %s <%s>\n",
        tag.Commit.Committer.Name, tag.Commit.Committer.Email)
    fmt.Printf("CommitDate:  %s\n",
        tag.Commit.CommitterDate.Format(time.RubyDate))
    fmt.Printf("Message:     %s\n\n", tag.Commit.Message)
    if tag.Commit.Signature != nil {
        fmt.Printf("Fingerprint: %s\n", tag.Commit.Signature.Fingerprint)
    }
}
```

```{ .text .no-select .no-copy }
Tagger:      purpleclay <purpleclaygh@gmail.com>
TaggerDate:  Thu Jun 29 07:05:18 +0100 2023
Message:     chore: tagged for release 0.3.2

Author:      Purple Clay <purpleclaygh@gmail.com>
AuthorDate:  Thu Jun 29 06:40:51 +0100 2023
Commit:      GitHub <noreply@github.com>
CommitDate:  Thu Jun 29 06:40:51 +0100 2023
Message:     fix: imported gpg key fails to sign when no tty is present (#33)
Fingerprint: 4AEE18**********
```

## Inspect a commit

Retrieve information about a specific commit by passing its reference to `ShowCommits`. Including its GPG signature, if present.

## Inspect a tree

Call `ShowTrees` to retrieve a listing of all files and directories within a specific tree index of a repository.

```{ .go .select linenums="1" }
package main

import (
    "fmt"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Query the gittest directory tree within gitz

    tree, _ := client.ShowTrees("ad4a68f6628ba9a6c367fe213eb8136fdb95ebcd")
    for _, entry := range tree[0].Entries {
        fmt.Printf("%s\n", entry)
    }
}
```

```{ .text .no-select .no-copy }
log.go
log_test.go
repository.go
repository_test.go
```

## Inspect a blob

Retrieve the contents of a file (blob) from a repository by passing its reference to `ShowBlobs`.
