---
icon: material/tag-outline
title: Creating, listing, and deleting tags
description: Manage tag creation, deletion, and retrieval within the current repository
---

# Creating, listing, and deleting tags

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-tag)

Manage tag creation, deletion, and retrieval within the current repository.

## Creating a tag

Tag a specific time point within a repository's history ready for pushing back to the configured remote. Tagging comes in two flavors, a `lightweight` and `annotated` tag. The main difference being an annotated tag is treated as a complete object within git and must include a message (_or annotation_). When querying the git history, an annotated tag will contain details such as the author and its GPG signature, if signed.

Gitz supports both tags but defaults to creating the lightweight variant unless instructed.

### Creating a lightweight tag

Calling `Tag` with a valid name[^1] will tag the repository with a lightweight tag:

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Tag("0.1.0")
    if err != nil {
        log.Fatal("failed to tag repository with version 0.1.0")
    }
}
```

??? question "Do you know the git tag naming restrictions?"

    1. Tags cannot begin, end with, or contain multiple consecutive `/` characters.
    1. Tags cannot contain any of the following characters: `\ ? ~ ^ : * [ @`
    1. Tags cannot contain a space ` `.
    1. Tags cannot end with a dot `.` or contain two consecutive dots `..` anywhere within them.

### Creating an annotated tag

Use the `WithAnnotation` option to switch to annotated tag creation mode:

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Tag("0.1.0", git.WithAnnotation("created tag 0.1.0"))
    if err != nil {
        log.Fatal("failed to tag repository with version 0.1.0")
    }
}
```

If you were to inspect the annotated tag, details about the author are now included:

```{ .text .no-select .no-copy }
$ git show 0.1.0

tag 0.1.0
Tagger: Purple Clay <**********(at)*******>
Date:   Mon Feb 20 05:58:55 2023 +0000

created tag 0.1.0

associated commit message
```

### Creating a local tag

Use the `WithLocalOnly` option to prevent a tag from being pushed back to the remote.

### Tagging a specific commit :material-new-box:{.new-feature title="Feature added on the 19th of September 2023"}

Use the `WithCommitRef` option to ensure a specific commit within the history is tagged.

### Batch tagging :material-new-box:{.new-feature title="Feature added on the 19th of September 2023"}

By calling `TagBatch`, a batch of tags can be created. `gitz` will enforce the `WithLocalOnly` option before pushing them to the remote in one transaction.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.TagBatch([]string{"1.0.0", "1.0", "1"})
    if err != nil {
        log.Fatal("failed to batch tag repository")
    }
}
```

A batch of tags that target specific commits can also be created by calling `TagBatchAt`. `gitz` will enforce the `WithLocalOnly` and `WithCommitRef` options before pushing them back to the remote in one transaction.

```{ .go .no-select linenums="1" }
client.TagBatchAt([]string{"0.1.0", "740a8b9", "0.2.0", "9e7dfbb"})
```

## Retrieving all tags

Calling `Tags` will retrieve all tags from the current repository in ascending lexicographic order:

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tags 0.9.0, 0.9.1, 0.10.0 and 0.11.0

    tags, err := client.Tags()
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

The resulting output would be:

```{ .text .no-select .no-copy }
0.10.0
0.11.0
0.9.0
0.9.1
```

### Changing the sort order

You can change the default sort order when retrieving tags by using the `WithSortBy` option. Various [sort keys](https://git-scm.com/docs/git-for-each-ref#_field_names) exist, each affecting the overall sort differently. If using multiple sort keys, the last one becomes the primary key. Prefix any key with a `-` for a descending sort. For convenience `gitz` provides constants for the most common sort keys:

- `creatordate`: sort by the creation date of the associated commit.
- `refname`: sort by the tags reference name in lexicographic order (_default_).
- `taggerdate`: sort by the tags creation date.
- `version:refname`: interpolates the tag as a version number and sorts.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tags 0.9.0, 0.9.1, 0.10.0 and 0.11.0

    tags, err := client.Tags(git.WithSortBy(git.VersionDesc))
    if err != nil {
        log.Fatal("failed to sort local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

The resulting output would now be:

```{ .text .no-select .no-copy }
0.11.0
0.10.0
0.9.1
0.9.0
```

### Filtering by pattern

Filter local tags using pattern-based git [shell globs](https://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm) with the `WithShellGlob` option. If using multiple patterns, a tag only needs to match one.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tags 0.8.0, 0.9.0, 1.0.0 and v1

    tags, err := client.Tags(git.WithShellGlob("*.*.*"))
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

The filtered output would be:

```{ .text .no-select .no-copy }
0.8.0
0.9.0
1.0.0
```

### User-defined filters

Extend filtering by applying user-defined filters to the list of retrieved tags with the `WithFilters` option. Execution of filters is in the order defined.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

var (
    uiFilter = func(tag string) bool {
        return strings.HasPrefix(tag, "ui/")
    }

    noVTagsFilter := func(tag string) bool {
        return !strings.HasSuffix(tag, "v1")
    }
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tags ui/1.0.0, ui/v1, backend/1.0.0, backend/v1

    tags, err := client.Tags(git.WithFilters(uiFilter, noVTagsFilter))
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

The filtered output would be:

```{ .text .no-select .no-copy }
ui/0.1.0
```

### Limiting the number of tags

Define the maximum number of returned tags through the `WithCount` option. Limiting is applied as a post-processing step after all other options.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tags 0.1.0, 0.2.0, 0.3.0, 0.4.0

    tags, err := client.Tags(git.WithSortBy(git.VersionDesc),
        git.WithCount(2))
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

The limited output would be:

```{ .text .no-select .no-copy }
0.4.0
0.3.0
```

## Deleting a tag

Call `DeleteTag` to delete a local tag and sync it with the remote:

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tag 0.1.0

    tags, err := client.DeleteTag("0.1.0")
    if err != nil {
        log.Fatal("failed to delete tag)
    }
}
```

[^1]: Gitz defers the validation of a tag name to the git client. Any error is captured and returned back to the caller

### Only delete local reference

To prevent a deletion from being pushed back to the remote, use the `WithLocalDelete` option.

## Deleting multiple tags

Call `DeleteTags` if you need to delete a batch of tags and sync it with the remote. Use the `WithLocalDelete` option to prevent any deletion from being pushed back to the remote.

## Signing a tag using GPG

Any tag against a repository can be GPG signed by the tagger to prove its authenticity through GPG verification. By setting the `tag.gpgSign` and `user.signingKey` git config options, GPG signing, can become an automatic process. `gitz` provides options to control this process and manually overwrite existing settings per tag.

### Annotating the signed tag

A signed tag must have an annotation. `gitz` defaults this to `created tag <ref>`, but you can change this with the `WithAnnotation` option.

### Sign an individual tag

If the `tag.gpgSign` git config setting is not enabled, you can selectively GPG sign a tag using the `WithSigned` option.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Tag("0.1.0", git.WithSigned())
    if err != nil {
        log.Fatal("failed to gpg sign tag with user.signingKey")
    }
}
```

### Select a GPG signing key

If multiple GPG keys exist, you can cherry-pick a key when tagging using the `WithSigningKey` option, overriding the `user.signingKey` git config setting, if set.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Tag("0.1.0", git.WithSigningKey("E5389A1079D5A52F"))
    if err != nil {
        log.Fatal("failed to gpg sign tag with provided public key")
    }
}
```

### Prevent a tag from being signed

You can disable the GPG signing of a tag by using the `WithSkipSigning` option, overriding the `tag.gpgSign` git config setting if set.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()
    client.Tag("0.1.0", git.WithSkipSigning())
}
```

## Providing git config at execution

You can provide git config through the `WithTagConfig` option to only take effect during the execution of a `Tag`, removing the need to change config permanently.
