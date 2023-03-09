---
icon: material/test-tube
status: new
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

Initialize a repository with a predefined log by using the `WithLog` option. It can contain both commit messages and lightweight tags and is written to the repository in reverse chronological order. The expected format is equivalent to the output from the git command: `git log --pretty='format:%d %s'`.

```{ .go .select linenums="1" }
func TestGreatFeature(t *testing.T) {
    log := `(tag: 0.1.0) feat: this is a brand new feature
docs: write amazing material mkdocs documentation
ci: include github release workflow`
    gittest.InitRepository(t, gittest.WithLog(log))

    // test logic and assertions to follow ...
}
```

### With a remote log

Initialize the remote origin of a repository with a predefined log using the `WithRemoteLog` option. Ideal for simulating a delta between the current log and its remote counterpart.

```{ .go .select linenums="1" }
func TestGreatFeature(t *testing.T) {
    log := `(tag: 0.1.0) feat: this is a brand new feature
docs: write amazing material mkdocs documentation
ci: include github release workflow`
    gittest.InitRepository(t, gittest.WithRemoteLog(log))

    // test logic and assertions to follow ...
}
```

### With untracked files

Create a set of untracked files within a repository using the `WithFiles` option. File paths can be fully qualified or relative to the repository root. Each created file will contain a sample of `lorem ipsum` text.

```{ .go .select linenums="1" }
func TestGreatFeature(t *testing.T) {
    gittest.InitRepository(t, gittest.WithFiles("a.txt", "dir/b.txt"))

    // test logic and assertions to follow ...
}
```

Querying the repository status:

```text
$ git status --porcelain
?? a.txt
?? dir/b.txt
```

### With staged files

Create a set of staged (or tracked) files within a repository using the `WithStagedFiles` option.

```{ .go .select linenums="1" }
func TestGreatFeature(t *testing.T) {
    gittest.InitRepository(t, gittest.WithStagedFiles("a.txt", "dir/b.txt"))

    // test logic and assertions to follow ...
}
```

Querying the repository status:

```text
$ git status --porcelain
A  a.txt
A  dir/b.txt
```

### With local commits

Generate a set of local empty commits, ready to be pushed back to the remote, with the `WithLocalCommits` option. Generated Commits will be in chronological order.

```{ .go .select linenums="1" }
func TestGreatFeature(t *testing.T) {
    commits := []string{
        "docs: my first local commit",
        "fix: my second local commit",
        "feat: my third local commit",
    }
    gittest.InitRepository(t, gittest.WithLocalCommits(commits...))

    // test logic and assertions to follow ...
}
```

### With clone depth :material-new-box:{.new-feature title="Feature added on the 10th March 2023"}

Shallow clone a repository by truncating its history to a set depth.

```{ .go .select linenums="1" }
func TestGreatFeature(t *testing.T) {
    log := `(tag: 0.1.0) feat: this is a brand new feature
docs: write amazing material mkdocs documentation
ci: include github release workflow`
    gittest.InitRepository(t,
        gittest.WithLog(log), gittest.WithCloneDepth(1))

    // test logic and assertions to follow ...
}
```

Querying the repository log:

```text
$ git log --pretty=oneline --no-decorate --no-color

e03726d3c24bbbab106bd1ac6231c030e1296eb9 feat: this is a brand new feature
```

### Option initialization order

You can use any combination of options during repository initialization, but a strict order is applied.

1. `WithLog`: log history imported, both local and remote are in sync.
1. `WithCloneDepth`: shallow clone at the required depth.
1. `WithRemoteLog`: remote log history imported, creating a delta between local and remote.
1. `WithLocalCommits`: local commits created and not pushed back to remote.
1. `WithFiles` and `WithStagedFiles`: files generated and staged if needed.
