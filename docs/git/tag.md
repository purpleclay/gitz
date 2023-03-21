---
icon: material/tag-outline
status: new
title: Tagging a repositories history
description: Tag a specific time point within a repository's history with a lightweight or annotated tag
---

# Tagging a repositories history

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-tag)

Tag a specific time point within a repository's history ready for pushing back to the configured remote. Tagging comes in two flavors, a `lightweight` and `annotated` tag. The main difference being an annotated tag is treated as a complete object within git and must include a message (_or annotation_). When querying the git history, an annotated tag will contain details such as the author and its GPG signature, if signed.

Gitz supports both tags but defaults to creating the lightweight variant unless instructed.

## Creating a Lightweight Tag

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

## Creating an Annotated Tag

Use the `WithAnnotation` option to switch to annotated tag creation mode:

```{ .go .select linenums="1" hl_lines="6" }
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

... # (1)!
```

1. The associated commit was replaced with a `...` for brevity

[^1]: Gitz defers the validation of a tag name to the git client. Any error is captured and returned back to the caller
