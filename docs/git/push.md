---
icon: material/arrow-right-bold-box-outline
---

# Pushing the latest changes to a Remote

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-push)

Push all local repository changes back to the remote, ensuring all references are tracked and both instances are in sync.

## Push committed changes back to the Remote

Calling `Push` will attempt to push all locally committed changes back to the remote for the current branch:

```{ .go .select linenums="1" }
package main

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

func main() {
    client, _ := git.NewClient()

    // tag 0.1.0 has been created and is tracked locally

    _, err := client.PushTag("0.1.0")
    if err != nil {
        log.Fatal("failed to push tag 0.1.0 to the remote")
    }
}
```
