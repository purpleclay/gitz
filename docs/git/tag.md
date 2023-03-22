---
icon: material/tag-outline
status: new
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

## Retrieving all tags :material-new-box:{.new-feature title="Feature added on the 21st March 2023"}

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

```text
0.10.0
0.11.0
0.9.0
0.9.1
```

### Changing the sort order :material-new-box:{.new-feature title="Feature added on the 21st March 2023"}

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

```text
0.11.0
0.10.0
0.9.1
0.9.0
```

### Filtering by pattern :material-new-box:{.new-feature title="Feature added on the 21st March 2023"}

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
        log.Fatal("failed to delete tag)
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

The filtered output would be:

```text
0.8.0
0.9.0
1.0.0
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
