---
icon: material/archive-lock-open-outline
title: Staging changes within a repository
description: Stage changes to a particular file or folder within the current repository for inclusion within the next commit
---

# Staging changes within a repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-stage)

Stage changes to a particular file or folder within the current repository for inclusion within the next commit. Staging is a prerequisite to committing and pushing changes back to the repository remote.

## Staging all changes

By default, all files (`tracked` and `untracked`) within the current repository are staged automatically unless explicitly ignored through a `.gitignore` file:

```{ .go .select linenums="1" }
package main

import (
	"log"

	git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // create multiple files within the following hierarchy:
    //  > a.txt
    //  > b.txt

    _, err := client.Stage()
    if err != nil {
        log.Fatal("failed to stage all files")
    }
}
```

And to verify the staged changes:

```{ .text .no-select .no-copy }
$ git status --porcelain

A  a.txt
A  b.txt
```

## Staging a file or folder

Cherry-picking the staging of files and folders is accomplished using the `WithPathSpecs` option:

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // create multiple files within the following hierarchy:
    //  > root.txt
    //  > folder
    //    > a.txt
    //    > b.txt

    _, err := client.Stage(git.WithPathSpecs("root.txt", "folder/a.txt"))
    if err != nil {
        log.Fatal("failed to stage files")
    }
}
```

And to verify the staged changes:

```{ .text .no-select .no-copy }
$ git status --porcelain

A  folder/a.txt
?? folder/b.txt
A  root.txt
```
