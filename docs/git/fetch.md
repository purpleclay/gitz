---
icon: material/clipboard-arrow-left-outline
status: new
title: Fetching the latest changes from a remote
description: Fetch all changes from a remote without integrating them into the current working directory
---

# Fetching the latest changes from a remote

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-fetch)

Fetch all remote changes from a remote repository without integrating (merging) them into the current repository (working directory). Ensures the existing repository only tracks the latest remote changes.

## Fetch all changes

Calling `Fetch` without any options will attempt to retrieve and track all the latest changes from the default remote.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Fetch()
    if err != nil {
        log.Fatal("failed to fetch all changes from the remote")
    }
}
```

## Fetch from all remotes

To fetch the latest changes from all tracked remotes, use the `WithAll` option.

## Fetch and follow tags

Retrieve all of the latest tags and track them locally with the `WithTags` option.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Existing locally tracked tags: 0.1.0
    // Additional tags that exist at the remote: 0.2.0, 0.3.0

    _, err := client.Fetch(git.WithTags())
    if err != nil {
        log.Fatal("failed to fetch all changes from the remote")
    }

    tags, err := client.Tags()
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

```{ .text .no-select .no-copy }
0.1.0
0.2.0
0.3.0
```

## Fetch but do not follow tags

The `WithIgnoreTags` option turns off local tracking of tags retrieved from the remote.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Existing locally tracked tags: 0.1.0
    // Additional tags that exist at the remote: 0.2.0, 0.3.0

    _, err := client.Fetch(git.WithIgnoreTags())
    if err != nil {
        log.Fatal("failed to fetch all changes from the remote")
    }

    tags, err := client.Tags()
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

```{ .text .no-select .no-copy }
0.1.0
```

## Limit fetching of commit history

Limit the number of commits fetched from the tip of each remote branch history, using the `WithDepthTo` option. This can be used to deepen or shorten the existing history of a shallow cloned repository.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Fetch(git.WithDepthTo(2))
    if err != nil {
        log.Fatal("failed to fetch all changes from the remote")
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
feat: add initial support for git fetch
feat: extend pull options to control how change sets are retrieved
```

## Force fetching into an existing local branch

Fetching may be refused if updating a locally tracked branch through the `WithFetchRefSpecs` option. Use the `WithForce` option to turn off this check.

## Providing git config at execution

You can provide git config through the `WithFetchConfig` option to only take effect during the execution of a `Fetch`, removing the need to change config permanently.
