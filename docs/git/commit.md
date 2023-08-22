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

## Allowing an empty commit

You can create empty commits without staging any files using the `WithAllowEmpty` option.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()
    client.Commit("no files are staged here", git.WithAllowEmpty())
}
```

## Signing a commit using GPG

Any commit to a repository can be GPG signed by an author to prove its authenticity through GPG verification. By setting the `commit.gpgSign` and `user.signingKey` git config options, GPG signing, can become an automatic process. `gitz` provides options to control this process and manually overwrite existing settings per commit.

### Sign an individual commit

If the `commit.gpgSign` git config setting is not enabled; you can selectively GPG sign a commit using the `WithGpgSign` option.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Commit("no files are staged here",
        git.WithAllowEmpty(),
        git.WithGpgSign())
    if err != nil {
        log.Fatal("failed to gpg sign commit with user.signingKey")
    }
}
```

### Select a GPG signing key

If multiple GPG keys exist, you can cherry-pick a key during a commit using the `WithGpgSigningKey` option, overriding the `user.signingKey` git config setting, if set.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Commit("no files are staged here",
        git.WithAllowEmpty(),
        git.WithGpgSigningKey("E5389A1079D5A52F"))
    if err != nil {
        log.Fatal("failed to gpg sign commit with provided public key")
    }
}
```

### Prevent a commit from being signed

You can disable the GPG signing of a commit by using the `WithNoGpgSign` option, overriding the `commit.gpgSign` git config setting, if set.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()
    client.Commit("prevent commit from being signed",
        git.WithAllowEmpty(),
        git.WithNoGpgSign())
}
```

## Providing git config at execution

You can provide git config through the `WithCommitConfig` option to only take effect during the execution of a `Commit`, removing the need to change config permanently.
