---
icon: material/clipboard-check-outline
status: new
---

# Checking the Integrity of a Repository

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
