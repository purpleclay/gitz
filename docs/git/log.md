---
icon: material/list-box-outline
---

# Inspecting the Commit Log of a Repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-log)

Retrieve the commit log of a repository in an easy-to-parse format.

## View the entire log

Calling `Log` without any options will retrieve the entire repository log from the current branch. A default formatting of `--pretty=oneline --no-decorate --no-color` is applied by `gitz` during log retrieval:

```{ .go .select linenums="1" }
package main

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

```{ .text .no-select }
99cb396148cff7db435cebb9a8ea95a11b5658e1 fix: parsing error when input string is too long
e8b67f90c613340fd61392fda03181d86a7febbe ci: extend the existing build workflow to include integration tests
c90e819baa6ddc212811924c51256e65f53d4c32 docs: create initial mkdocs material documentation
1635f10b81a810833b163793e6ef902d52a89789 feat: add first feature to library
a09348464773e99dbc94a5494b5b83b253c18019 initialized repository
```

By default, `gitz` parses the log into a structured output accessible through the `Commits` property. This structure contains each commit's associated `Hash`, `Abbreviated Hash`, and `Message`.

### Return only raw output from the log

Use the `WithRawOnly` option to skip parsing of the log into the structured `Commits` property, improving performance. Perfect if you want to carry out any custom processing.

```{ .go .select linenums="1" }
package main

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

```text
[]
```

## View the log from a point in time

The `WithRef` option provides a starting point other than HEAD (_most recent commit_) when retrieving the log history. A reference can be a `Commit Hash`, `Branch Name`, or `Tag`. Output from this option will be a shorter, fine-tuned log.

```{ .go .select linenums="1" }
package main

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

```text
d611a22c1a009bd74bc2c691b331b9df38828dae fix: typos in file content
9b342465255d1a8ec4f5517eef6049e5bcc8fb45 feat: a brand new feature
```
