---
icon: material/sheep
status: new
title: Cloning a repository
description: Clone a repository by its provided URL into a newly created directory
---

# Cloning a repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-clone)

Clone a repository by its provided URL into a newly created directory and check out an initial branch forked from the cloned repository’s default branch.

## Clone everything

Calling `Clone` will result in a repository containing all remote-tracking branches, tags, and a complete history.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Clone("https://github.com/purpleclay/gitz")
    if err != nil {
        log.Fatal("failed to clone repository")
    }
}
```

## With a truncated history

If a complete git history isn't desirable, a faster approach to cloning is to provide a clone depth using the `WithDepth` option. This results in a truncated history, better known as a `shallow clone`. Great for large established repositories.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // A repository exists with the following commits:
    // > feat: add support for shallow cloning
    // > initialized repository

    _, err := client.Clone("https://github.com/purpleclay/gitz",
        git.WithDepth(1))
    if err != nil {
        log.Fatal("failed to clone repository")
    }

    repoLog, err := client.Log()
    if err != nil {
        log.Fatal("failed to retrieve repository log")
    }

    for _, commit := range repoLog.Commits {
        fmt.Println(commit.Message)
    }
}
```

Printing the log results in:

```{ .text .no-select .no-copy }
feat: add options that can configure the size of the repository during a clone
```

## Without tags

To further streamline a clone, the `WithNoTags` option prevents the downloading of tags from the remote.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Clone("https://github.com/purpleclay/gitz",
        git.WithNoTags())
    if err != nil {
        log.Fatal("failed to clone repository")
    }

    tags, err := client.Tags()
    if err != nil {
        log.Fatal("failed to retrieve repository tags")
    }

    if len(tags) == 0 {
        fmt.Println("repository has no tags")
    }
}
```

Printing the output results in:

```{ .text .no-select .no-copy }
repository has no tags
```

## Clone into a named directory

By default, git will always clone the repository into a directory based on the human-friendly part of the clone URL. To clone into a different directory, use the `WithDirectory` option.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Clone("https://github.com/purpleclay/gitz",
        git.WithDirectory("my-gitz"))
    if err != nil {
        log.Fatal("failed to clone repository")
    }
}
```

## Clone a branch or tag

Git will always clone and checkout the default branch of a repository. To change this behavior, provide a branch or tag reference with the `WithCheckoutRef` option. The latter results in a `detached HEAD` where the `HEAD` of a repository points to a specific commit reference.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Clone("https://github.com/purpleclay/gitz",
        git.WithCheckoutRef("0.1.0"))
    if err != nil {
        log.Fatal("failed to clone repository")
    }
}
```
