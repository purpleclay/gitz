---
icon: material/keyboard-settings-outline
title: Managing your git config
description: Get and set your local git repository config
---

# Managing your git config

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-config)

Manage settings within your local git config, changing the behavior of the git client.

## Retrieve all settings

Retrieve all git config for the current repository using `Config`.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    cfg, err := client.Config()
    if err != nil {
        log.Fatal("failed to retrieve config for current repository")
    }

    fmt.Println(cfg["user.name"])
}
```

The value for the config setting would be:

```{ .text .no-select .no-copy }
purpleclay
```

## Retrieve a batch of settings

A batch of settings can be retrieved using `ConfigL` (_local_), `ConfigS` (_system_), or `ConfigG` (_global_). A partial retrieval is not supported and will fail if any are missing. All values for a setting are retrieved and ordered by the latest.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()
    cfg, err := client.ConfigL("user.name", "user.email")
    if err != nil {
        log.Fatal("failed to retrieve local config settings")
    }

    fmt.Println(cfg["user.name"][0])
    fmt.Println(cfg["user.email"][0])
}
```

The value for each config setting would be:

```{ .text .no-select .no-copy }
purpleclay
**********************
```

## Update a batch of settings

You can update multiple settings in a batch using `ConfigSetL` (_local_), `ConfigSetS` (_system_), or `ConfigSetG` (_global_). Pre-validation of config paths improves the chance of a successful update, but a partial batch may occur upon failure.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    err := client.ConfigSetL("custom.setting1", "value",
        "custom.setting2", "value")
    if err != nil {
        log.Fatal("failed to set local config settings")
    }
}
```
