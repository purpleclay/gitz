---
icon: material/clipboard-check-outline
title: Git checks and how to use them
description: A series of inbuilt check for inspecting the environment and current repository
---

# Git checks and how to use them

`gitz` comes with a series of inbuilt checks for inspecting the environment and current repository.

## Checking for the existence of a Git Client

When creating a new client, `gitz` will check for the existence of git using the `PATH` environment variable. An error is returned if no client exists.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, err := git.NewClient()
    if err != nil {
        log.Fatal(err.Error())
    }

    fmt.Println(client.Version())
}
```

## Checking the integrity of a Repository

Check the integrity of a repository by running a series of tests and capturing the results for inspection.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    repo, err := client.Repository()
    if err != nil {
        log.Fatal("failed to check the current repository")
    }

    fmt.Printf("Default Branch: %s\n", repo.DefaultBranch)
    fmt.Printf("Detached Head:  %t\n", repo.DetachedHead)
    fmt.Printf("Origin:         %s\n", repo.Origin)
    fmt.Printf("Remotes:        %#v\n", repo.Remotes)
    fmt.Printf("Root Directory: %s\n", repo.RootDir)
    fmt.Printf("Shallow Clone:  %t\n", repo.ShallowClone)
}
```

Example output when checking the integrity of a repository cloned within a CI system:

```{ .text .no-select .no-copy }
Default Branch: main
Detached Head:  true
Origin:         git@github.com:purpleclay/gitz.git
Remotes:        map[string]string{"origin":"git@github.com:purpleclay/gitz.git"}
Root Directory: /dev/github.com/purpleclay/gitz
Shallow Clone:  false
```
