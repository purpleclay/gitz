---
icon: material/keyboard-settings-outline
status: new
title: Managing your git config
description: Get and set your local git repository config
---

# Managing your git config

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-config) | :material-beaker-outline: Experimental

TODO

## Retrieve a local setting :material-new-box:{.new-feature title="Feature added on the 31st March of 2023"}

Providing a valid path to `Config` will retrieve all values associated with a setting in modification order.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()
    // setting user.name to cover up real identity

    cfg, err := client.Config("user.name")
    if err != nil {
        log.Fatal("failed to retrieve config setting")
    }

    for _, v := range cfg {
        fmt.Println(v)
    }
}
```

The values for the config setting would be:

```text
purpleclay
****
```

## Retrieve a batch of local settings :material-new-box:{.new-feature title="Feature added on the 31st March of 2023"}

TODO

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tags 0.1.0, 0.2.0, 0.3.0, 0.4.0

    tags, err := client.Tags(git.WithSortBy(git.VersionDesc),
        git.WithCount(2))
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

## Update a local setting :material-new-box:{.new-feature title="Feature added on the 31st March of 2023"}

TODO

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tags 0.1.0, 0.2.0, 0.3.0, 0.4.0

    tags, err := client.Tags(git.WithSortBy(git.VersionDesc),
        git.WithCount(2))
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```

## Updating a batch of local settings :material-new-box:{.new-feature title="Feature added on the 31st March of 2023"}

TODO

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // Repository contains tags 0.1.0, 0.2.0, 0.3.0, 0.4.0

    tags, err := client.Tags(git.WithSortBy(git.VersionDesc),
        git.WithCount(2))
    if err != nil {
        log.Fatal("failed to retrieve local repository tags")
    }

    for _, tag := range tags {
        fmt.Println(tag)
    }
}
```
