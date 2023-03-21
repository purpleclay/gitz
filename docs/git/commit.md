---
icon: material/archive-lock-outline
title: Committing changes to a repository
description: Create a commit within the current repository and describe those changes with a given log message
---

# Committing changes to a repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-commit)

Create a commit (_snapshot of changes_) within the current repository and describe those changes with a given log message. A commit will only exist within the local history until pushed back to the repository remote.

## Commit a Snapshot of Repository Changes

Calling `Commit` with a message will create a new commit within the repository:

```{ .go .select linenums="1" }
package main

import (
	"log"

	git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // stage all changes to files and folders

    _, err := client.Commit("feat: a brand new feature")
    if err != nil {
        log.Fatal("failed to commit latest changes within repository")
    }
}
```

And to verify its creation:

```{ .text .no-select .no-copy }
$ git log -n1

commit 703a6c9bc9ee91d0c226b169b131670fb92d9a0a (HEAD -> main)
Author: Purple Clay <**********(at)*******>
Date:   Mon Feb 20 20:43:49 2023 +0000

    feat: a brand new feature
```
