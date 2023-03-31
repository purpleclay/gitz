---
icon: material/keyboard-settings-outline
status: new
title: Managing your git config
description: Get and set your local git repository config
---

# Managing your git config

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-config) | :material-beaker-outline: Experimental

Manage settings within your local git config, changing the behavior of the git client.

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
    // setting user.name to purpleclay to cover up real identity

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

```{ .text .no-select .no-copy }
purpleclay
****
```

## Retrieve a batch of local settings :material-new-box:{.new-feature title="Feature added on the 31st March of 2023"}

For convenience, multiple local settings can be retrieved in a batch using `ConfigL`. A partial batch is not supported and will fail if any setting does not exist.

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
        log.Fatal("failed to retrieve config settings")
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

## Update a local setting :material-new-box:{.new-feature title="Feature added on the 31st March of 2023"}

To update a local git setting, call `ConfigSet`` with a path and corresponding value.

```{ .go .select linenums="1" }
package main

import (
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()
    err := client.ConfigSet("custom.setting", "value")
    if err != nil {
        log.Fatal("failed to set config setting")
    }
}
```

## Updating a batch of local settings :material-new-box:{.new-feature title="Feature added on the 31st March of 2023"}

Multiple local settings can be updated in a batch using `ConfigSetL`. Pre-validation of config paths improves the chance of a successful update, but a partial batch may occur upon failure.

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
        log.Fatal("failed to set config setting")
    }
}
```
