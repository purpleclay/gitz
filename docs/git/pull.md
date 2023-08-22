---
icon: material/arrow-left-bold-box-outline
status: new
title: Pulling the latest changes from a remote
description: Pull all changes from a remote and integrate them into the current working directory
---

# Pulling the latest changes from a remote

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-pull)

Pull all changes from a remote into the current working directory. Ensures the existing repository keeps track of changes and stays in sync.

## Pull the latest changes from the current Branch

Calling `Pull` will attempt to sync the current branch with its counterpart from the remote:

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // a new file was added to the hierarchy at the remote:
    //  > folder
    //    > c.txt

    out, err := client.Pull()
    if err != nil {
        log.Fatal("failed to pull latest changes from remote")
    }

    fmt.Println(out)
}
```

Printing the output from this command:

```{ .text .no-select .no-copy }
remote: Enumerating objects: 5, done.
remote: Counting objects: 100% (5/5), done.
remote: Compressing objects: 100% (3/3), done.
remote: Total 3 (delta 0), reused 0 (delta 0), pack-reused 0
Unpacking objects: 100% (3/3), 300 bytes | 150.00 KiB/s, done.
From /Users/paulthomas/dev/./gitrepo
   703a6c9..8e87f78  main       -> origin/main
Updating 703a6c9..8e87f78
Fast-forward
 folder/c.txt | 1 +
 1 file changed, 1 insertion(+)
 create mode 100644 folder/c.txt
```

## Configuring fetch behavior :material-new-box:{.new-feature title="Feature added on the 21st of August 2023"}

When pulling changes from a remote repository, a fetch happens before merging any changes. Configuring the behavior of this fetch is possible using the supported `WithFetch...` [options](./fetch.md).

For example, you can limit the fetched commit history with the `WithFetchDepthTo` option.

## Providing git config at execution

You can provide git config through the `WithPullConfig` option to only take effect during the execution of a `Pull`, removing the need to change config permanently.
