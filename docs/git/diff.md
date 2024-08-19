---
icon: material/vector-difference
title: Diff local changes within a repository
description: Inspect the current repository for changes to files and folders
status: new
---

#  Diff local changes within a repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-diff)

Show unified changes to local files within the current repository.

## Diff all changes

Calling `Diff` without options will retrieve all changes within the current repository.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Changes are made to local files

    _, err := client.Diff()
    if err != nil {
        log.Fatal("failed to diff repository for changes")
    }
}
```

##  Diff changes for specific files or folders

To only retrieve changes for specific files or folders, use the `WithDiffPaths` option.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Changes are made to local files

    _, err := client.Diff(git.WithDiffPaths("main.go", "internal/cache"))
    if err != nil {
        log.Fatal("failed to diff repository for changes")
    }
}
```
