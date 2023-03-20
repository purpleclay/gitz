---
icon: material/sheep
status: new
---

# Cloning a Repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-clone)

Clone a repository by its provided URL into a newly created directory and check out an initial branch forked from the cloned repository’s default branch.

## Clone everything :material-new-box:{.new-feature title="Feature added on the 10th March 2023"}

Calling `Clone` will result in a repository containing all remote-tracking branches, tags, and a complete history.

```{ .go .select linenums="1" }
package main

import (
	"log"

	git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    _, err := client.Clone("https://github.com/purpleclay/gitz")
    if err != nil {
        log.Fatal("failed to clone repository")
    }
}
```
