---
icon: material/list-box-outline
title: Inspecting the commit log of a repository
description: Retrieve the commit log of a repository in an easy-to-parse format
---

# Inspecting the commit log of a repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-log)

Retrieve the commit log of a repository in an easy-to-parse format.

## View the entire log

Calling `Log` without any options will retrieve the entire repository log from the current branch. A default formatting of `--pretty=oneline --no-decorate --no-color` is applied by `gitz` during log retrieval:

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    log, err := client.Log()
    if err != nil {
        log.Fatal("failed to retrieve repository log")
    }

    fmt.Println(log.Raw)
}
```

Printing the `Raw` output from this command:

```{ .text .no-select .no-copy }
99cb396148cff7db435cebb9a8ea95a11b5658e1 fix: input string parsing error
e8b67f90c613340fd61392fda03181d86a7febbe ci: extend existing build workflow
c90e819baa6ddc212811924c51256e65f53d4c32 docs: create mkdocs documentation
1635f10b81a810833b163793e6ef902d52a89789 feat: add first feature to library
a09348464773e99dbc94a5494b5b83b253c18019 initialized repository
```

By default, `gitz` parses the log into a structured output accessible through the `Commits` property. This structure contains each commit's associated `Hash`, `Abbreviated Hash`, and `Message`.

### Return only raw output from the log

Use the `WithRawOnly` option to skip parsing of the log into the structured `Commits` property, improving performance. Perfect if you want to carry out any custom processing.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    log, err := client.Log(git.WithRawOnly())
    if err != nil {
        log.Fatal("failed to retrieve repository log")
    }

    fmt.Println(log.Commits)
}
```

Printing the `Commits` property should now be an empty slice:

```{ .text .no-select .no-copy }
[]
```

## View the log from a point in time

When retrieving the log history, the `WithRef` option provides a starting point other than HEAD (_most recent commit_). A reference can be a `Commit Hash`, `Branch Name`, or `Tag`. Output from this option will be a shorter, fine-tuned log.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    log, err := client.Log(git.WithRef("0.1.0"))
    if err != nil {
        log.Fatal("failed to retrieve log for tag 0.1.0")
    }

    fmt.Println(log.Raw)
}
```

## View a snapshot of the log

The `WithRefRange` option provides a start and end point for retrieving a snapshot of the log history between two points in time. A reference can be a `Commit Hash`, `Branch Name`, or `Tag`. Output from this option will be a shorter, fine-tuned log.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    log, err := client.Log(git.WithRefRange("0.2.0", "0.1.0"))
    if err != nil {
        log.Fatal("failed to retrieve log between tags 0.2.0 and 0.1.0")
    }

    fmt.Println(log.Raw)
}
```

## View the log of any files or folders

Fine-tune the log history further with the `WithPaths` option. Providing a set of relative paths to any files and folders within the repository will include only commits related to their history.

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // commits are made to the following files in a newly initialized
    // repository:
    //  > a.txt
    //    ~ fix: typos in file content
    //  > b.txt
    //    ~ chore: restructure file into expected format
    //  > dir1 (b.txt, c.txt)
    //    ~ feat: a brand new feature

    log, err := client.Log(git.WithPaths("a.txt", "dir1"))
    if err != nil {
        log.Fatal("failed to retrieve log using custom paths")
    }

    fmt.Println(log.Raw)
}
```

Printing the `Raw` output from this command:

```{ .text .no-select .no-copy }
d611a22c1a009bd74bc2c691b331b9df38828dae fix: typos in file content
9b342465255d1a8ec4f5517eef6049e5bcc8fb45 feat: a brand new feature
```

## Cherry-picking a section of the log

Cherry-pick a section of the log by skipping and taking a set number of entries using the respective `WithSkip` and `WithTake` options. If combined, skipping has a higher order of precedence:

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // assuming the repository has the current history:
    //  ~ docs: document fix
    //  ~ fix: filtering on unsupported prefixes
    //  ~ docs: create docs using material mkdocs
    //  ~ feat: add new support for filtering based on prefixes
    //  ~ initialized the repository

    log, err := client.Log(git.WithSkip(1), git.WithTake(2))
    if err != nil {
        log.Fatal("failed to retrieve log using custom paths")
    }

    fmt.Println(log.Raw)
}
```

Printing the `Raw` output from this command:

```{ .text .no-select .no-copy }
9967e3c6196422a6a97afa4b6fca9f609bb5490b fix: filtering on unsupported prefixes
1b1f4a725cfe44d5c9bd992be59f1130ed9d9911 docs: create docs using material mkdocs
```

## Filtering the log with pattern matching

Filter the commit log to only contain entries that match any set of patterns (_regular expressions_) using the `WithGrep` option:

```{ .go .select linenums="1" }
package main

import (
    "fmt"
    "log"

    git "github.com/purpleclay/gitz"
)

func main() {
    client, _ := git.NewClient()

    // assuming the repository has the current history:
    //  ~ fix: forgot to trim whitespace from patterns
    //  ~ docs: document pattern matching option
    //  ~ feat: filter log with pattern matching
    //  ~ initialized the repository

    log, err := client.Log(git.WithGrep("^docs", "matching$"))
    if err != nil {
        log.Fatal("failed to retrieve log with pattern matching")
    }

    fmt.Println(log.Raw)
}
```

Printing the `Raw` output from this command, with matches highlighted for reference only:

```{ .text .no-select .no-copy }
2d68a506fe7d5148db0a10ea143752991a65c26d {==docs==}: document pattern matching option
5bfd532328ed2e9ea6d3062eb3a331f42468a7e3 feat: filter log with pattern {==matching==}
```

### Filter entries that do not match

Combining the `WithInvertGrep` and `WithGrep` options will inverse pattern matching and filter on log entries that do not contain any of the provided patterns.

### Filter entries that match all patterns

Pattern matching uses `or` semantics by default, matching on log entries that satisfy any of the defined patterns. You can change this behavior to match against all patterns using `and` semantics with the `WithMatchAll` option.
