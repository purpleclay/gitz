---
icon: material/clipboard-check-outline
status: new
---

# Git checks and how to use them

`gitz` comes with a series of inbuilt checks for inspecting the environment and current repository.

## Checking for the existence of a Git Client

When creating a new client, `gitz` will check for the existence of git using the `PATH` environment variable. An error is returned if no client exists.

```{ .go .select linenums="1" }
package main

func main() {
    client, err := git.NewClient()
    if err != nil {
        log.Fatal(err.Error())
    }

    fmt.Println(client.Version())
}
```

## Checking the integrity of a Repository :material-new-box:{.new-feature title="Feature added on the 10th March 2023"}

Check the integrity of a repository by running a series of tests and capturing the results for inspection.

```{ .go .select linenums="1" }
package main

func main() {
    client, _ := git.NewClient()

    repo, err := client.Repository()
    if err != nil {
        log.Fatal("failed to check the current repository")
    }

    fmt.Printf("Default Branch: %s\n", repo.DefaultBranch)
    fmt.Printf("Shallow Clone:  %t\n", repo.ShallowClone)
    fmt.Printf("Detached Head:  %t\n", repo.DetachedHead)
}
```

Example output when checking the integrity of a repository cloned within a CI system:

```{ .text .no-select }
Default Branch: main
Shallow Clone:  false
Detached Head:  true
```
