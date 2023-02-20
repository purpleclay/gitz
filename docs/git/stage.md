---
icon: material/archive-lock-open-outline
---

# Staging changes within a Repository

[:simple-git:{ .git-icon } Git Documentation](https://git-scm.com/docs/git-stage)

Stage changes to a particular file or folder within the current repository for inclusion within the next commit. Staging is a prerequisite to committing and pushing changes back to the repository remote.

## Staging a File or Folder

Calling `Stage` with a relative path to an individual file or folder will stage any changes:

```{ .go .select linenums="1" }
package main

func main() {
    client, _ := git.NewClient()

    // create multiple files within the following hierarchy:
    //  > root.txt
    //  > folder
    //    > a.txt
    //    > b.txt

    _, err := client.Stage("root.txt")
    if err != nil {
        log.Fatal("failed to stage file root.txt")
    }

    _, err := client.Stage("folder/")
    if err != nil {
        log.Fatal("failed to stage all changes within directory folder/")
    }
}
```
