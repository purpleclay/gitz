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

	verif, err := gitc.VerifyCommit(gLog.Commits[0].Hash)
	if err != nil {
		log.Fatal(err.Error())
	}

	if verif.Author.Name != "batman" {
		log.Fatalf("invalid author name, expecting: 'batman' but received: '%s'", verif.Author.Name)
	}

	if verif.Author.Email != "batman@dc.com" {
		log.Fatalf("invalid author email, expecting: 'batman@dc.com' but received: '%s'", verif.Author.Email)
	}

	if verif.Committer.Name != "batman" {
		log.Fatalf("invalid committer name, expecting: 'batman' but received: '%s'", verif.Committer.Name)
	}

	if verif.Committer.Email != "batman@dc.com" {
		log.Fatalf("invalid committer email, expecting: 'batman@dc.com' but received: '%s'", verif.Committer.Email)
	}

	if verif.Message != "this is a signed commit" {
		log.Fatalf("invalid commit message, expecting: 'this is a signed commit' but received: '%s'", verif.Message)
	}

	if verif.Signature.Fingerprint != gpgFingerprint {
		log.Fatalf("invalid fingerprint, expecting: '%s' but received: '%s'", gpgFingerprint, verif.Signature.Fingerprint)
	}

	if verif.Signature.Author.Name != "batman" {
		log.Fatalf("invalid signed-by name, expecting: 'batman' but received: '%s'", verif.Signature.Author.Name)
	}

	if verif.Signature.Author.Email != "batman@dc.com" {
		log.Fatalf("invalid signed-by email, expecting: 'batman@dc.com' but received: '%s'", verif.Signature.Author.Email)
	}
}
