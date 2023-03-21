---
icon: material/arrow-right-bold-box-outline
status: new
title: Pushing the latest changes back to a remote
description: Push all local repository changes back to the remote
---

# Pushing the latest changes back to a remote

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-push)

Push all local repository changes back to the remote, ensuring the remote now tracks all references.

## Push committed changes back to the Remote

Calling `Push` will attempt to push all locally committed changes back to the remote for the current branch:

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

## Push the created tag back to the Remote

Calling `PushTag` will attempt to push the newly created tag back to the remote:

```{ .go .select linenums="1" }
package main

import (
	"log"

	git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // tag 0.1.0 has been created and is tracked locally

    _, err := client.PushTag("0.1.0")
    if err != nil {
        log.Fatal("failed to push tag 0.1.0 to the remote")
    }
}
```
