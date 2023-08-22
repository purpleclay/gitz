---
icon: material/arrow-right-bold-box-outline
title: Pushing the latest changes back to a remote
description: Push all local repository changes back to the remote
---

# Pushing the latest changes back to a remote

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-push)

Push all local repository changes back to the remote, ensuring the remote now tracks all references.

## Pushing locally committed changes

Calling `Push` without any options will attempt to push all locally committed changes back to the remote for the current branch:

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // all changes have been staged and committed locally

    _, err := client.Push()
    if err != nil {
        log.Fatal("failed to push committed changes to the remote")
    }
}
```

## Pushing all local branches

To push changes spread across multiple branches back to the remote in a single atomic operation, use the `WillAllBranches` option:

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // modifications are made to multiple files across two
    // different branches
    //
    // b: new-feature
    //  > client.go
    // b: new-bug-fix
    //  > parser.go

    _, err := client.Push(git.WithAllBranches())
    if err != nil {
        log.Fatal("failed to stage files")
    }
}
```

## Pushing all local tags

All locally created tags can also be pushed back to the remote in a single atomic operation using the `WithAllTags` option:

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // multiple tags are created locally, 1.0.0 and v1

    _, err := client.Push(git.WithAllTags())
    if err != nil {
        log.Fatal("failed to stage files")
    }
}
```

## Cherry-pick what is pushed to the remote

The `WithRefSpecs` option provides greater freedom to cherry-pick locally created references (_branches and tags_) and push them back to the remote. A reference can be as simple as a name or as explicit as providing a source (_local_) to destination (_remote_) mapping. Please read the official git specification on how to construct [refspecs](https://git-scm.com/docs/git-push#Documentation/git-push.txt-ltrefspecgt82308203).

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // new branch and tag are created locally

    _, err := client.Push(git.WithRefSpecs("0.1.0", "new-branch"))
    if err != nil {
        log.Fatal("failed to stage files")
    }
}
```

## Push options

Support the transmission of arbitrary strings to the remote server using the `WithPushOptions` option.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // all changes have been staged and committed locally

    _, err := client.Push(git.WithPushOptions("ci.skip=true"))
    if err != nil {
        log.Fatal("failed to push committed changes to the remote")
    }
}
```

## Deleting references from the remote

Delete any number of references from the remote by using the `WithDeleteRefSpecs` option.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // a tag and branch have been deleted locally

    _, err := client.Push(git.WithDeleteRefSpecs("branch", "0.1.0"))
    if err != nil {
        log.Fatal("failed to delete references from the remote")
    }
}
```

## Providing git config at execution

You can provide git config through the `WithPushConfig` option to only take effect during the execution of a `Push`, removing the need to change config permanently.
