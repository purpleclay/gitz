---
icon: material/source-branch-sync
title: Checking out a branch
description: Clone a repository by its provided URL into a newly created directory
---

# Checking out a branch

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-checkout)

Switch from the default branch of a repository (working directory) to a new or existing one, syncing any file changes. All future actions are associated with this branch.

## Context-aware checking out :material-new-box:{.new-feature title="Feature added on the 23rd March of 2023"}

During a checkout, `gitz` inspects the repository for the existence of a branch and intelligently switches between creating a new one or checking out the existing reference.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    out, err := client.Checkout("a-new-branch")
    if err != nil {
        log.Fatal("failed to checkout branch")
    }

    fmt.Println(out)
}
```

If you were to print the output from the command, you would see a branch creation:

```{ .text .no-select .no-copy }
Switched to a new branch 'a-new-branch'
```

If you check out a branch that already exists, you will see a different output:

```{ .text .no-select .no-copy }
Switched to branch 'existing-branch'
Your branch is up to date with 'origin/existing-branch'.
```
