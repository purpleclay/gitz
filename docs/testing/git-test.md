---
icon: material/test-tube
---

# Testing your Interactions with Git

`gitz` includes a `gittest` package that enables interactions with Git to be unit tested within your projects. Add the following import to any of your test files to get started:

```{ .go .no-select }
import "github.com/purpleclay/gitz/gittest"
```

## Building a test repository

:octicons-beaker-24: Experimental

Only a single line of code is needed to initialize a test repository. And don't worry; it gets deleted after test execution.

```{ .go .select linenums="1" hl_lines="10" }
package main_test

import (
    "testing"

    "github.com/purpleclay/gitz/gittest"
)

func TestGreatFeature(t *testing.T) {
    gittest.InitRepository(t)

    // test logic and assertions to follow ...
}
```

Where `gittest` shines is in its ability to customize a repository during initialization through a set of options.

### With a commit log

`WithLog`

```{ .go .select linenums="1" }

```

### With a remote log

`WithRemoteLog`

```{ .go .select linenums="1" }

```

### With untracked files

`WithFiles`

```{ .go .select linenums="1" }

```

### With staged files

`WithStagedFiles`

```{ .go .select linenums="1" }

```

### With local commits

`WithLocalCommits`

```{ .go .select linenums="1" }

```
