package git_test

import (
	"fmt"
	"log"

	git "github.com/purpleclay/gitz"
)

func a() {
	client, err := git.NewClient()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(client.Version())
}
