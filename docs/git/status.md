---
icon: material/format-list-checkbox
status: new
title: Inspecting the status of a repository
description: Check the status of the current repository and identify if any changes exist
---

# Inspecting the status of a repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-status)

Indentify if any differences exist between the git staging area (_known as the index_) and the latest commit.

## Porcelain status :material-new-box:{.new-feature title="Feature added on the 25th of July 2023"}

To retrieve a parseable list of a changes within a repository, call `PorcelainStatus`. Changes are listed using the porcelain V1 format, consisting of a two character indicator, followed by a path to the identified change.

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

A two character indicator `' A'` denotes the status of a file. It should be read as the status within the index, followed by the status within the working tree. Staging a file will move it from the working tree to the index, moving the indicator from right (`' A'`) to left (`'A '`).

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

## Check if a repository is clean :material-new-box:{.new-feature title="Feature added on the 25th of July 2023"}

Calling `Clean` will return `true` if a repository has no oustanding changes.

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
