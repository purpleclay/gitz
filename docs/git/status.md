---
icon: material/format-list-checkbox
title: Inspecting the status of a repository
description: Check the status of the current repository and identify if any changes exist
---

# Inspecting the status of a repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-status)

Identify if any differences exist between the git staging area (_known as the index_) and the latest commit.

## Porcelain status

To retrieve a parseable list of changes within a repository, call `PorcelainStatus`. Changes are listed using the porcelain V1 format, consisting of a two-character indicator followed by a path to the identified change.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Add some files:
    //  new.txt
    //  staged.txt
    status, _ := client.PorcelainStatus()

    for _, s := range status {
        fmt.Printf("%s\n", s)
    }
}
```

```{ .text .no-select .no-copy }
?? new.txt
A  staged.txt
```

### Supported indicators

A two-character indicator, `' A'`, denotes the status of a file. It should be read as its status within the index, followed by its status within the working tree. Staging a file will move it from the working tree to the index, moving the indicator from the right (`' A'`) to the left (`'A '`).

```{ .text .no-select .no-copy }
'A' Added
'C' Copied
'D' Deleted
'!' Ignored
'M' Modified
'R' Renamed
'T' Type Changed (e.g. regular file to symlink)
'U' Updated
' ' Unmodified
'?' Untracked
```

## Check if a repository is clean

Calling `Clean` will return `true` if a repository has no outstanding changes.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()
    clean, _ := client.Clean()

    fmt.Printf("Is Clean: %t\n", clean)
}
```

```{ .text .no-select .no-copy }
Is Clean: true
```
