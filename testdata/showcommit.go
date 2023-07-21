package main

import (
	"log"
	"os"

	git "github.com/purpleclay/gitz"
)

var (
	gpgPublicKeyID = os.Getenv("GPG_PUBLIC_KEY_ID")
	gpgFingerprint = os.Getenv("GPG_FINGERPRINT")
)

func main() {
	gitc, _ := git.NewClient()
	gitc.ConfigSetL("user.signingkey", gpgPublicKeyID)

	if _, err := gitc.Commit("this is a signed commit", git.WithGpgSign(), git.WithAllowEmpty()); err != nil {
		log.Fatal(err.Error())
	}
	gLog, _ := gitc.Log(git.WithTake(1))

	commits, err := gitc.ShowCommits(gLog.Commits[0].Hash)
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(commits) != 1 {
		log.Fatalf("invalid number of commits, expected '1' commit but recevied: '%d'", len(commits))
	}

	commit := commits[gLog.Commits[0].Hash]

	if commit.Signature.Fingerprint != gpgFingerprint {
		log.Fatalf("invalid fingerprint, expecting: '%s' but received: '%s'", gpgFingerprint, commit.Signature.Fingerprint)
	}

	if commit.Signature.Author.Name != "batman" {
		log.Fatalf("invalid signed-by name, expecting: 'batman' but received: '%s'", commit.Signature.Author.Name)
	}

	if commit.Signature.Author.Email != "batman@dc.com" {
		log.Fatalf("invalid signed-by email, expecting: 'batman@dc.com' but received: '%s'", commit.Signature.Author.Email)
	}
}
