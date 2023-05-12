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
	gitc, err := git.NewClient()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := gitc.ConfigSetL("user.signingkey", gpgPublicKeyID); err != nil {
		log.Fatal(err.Error())
	}

	if _, err := gitc.Tag("0.1.0", git.WithSigned(), git.WithAnnotation("withsigned")); err != nil {
		log.Fatal(err.Error())
	}

	tag, err := gitc.VerifyTag("0.1.0")
	if err != nil {
		log.Fatal(err.Error())
	}

	if tag.Tagger.Name != "batman" {
		log.Fatalf("invalid tagger name, expecting: 'batman' but received: '%s'", tag.Tagger.Name)
	}

	if tag.Tagger.Email != "batman@dc.com" {
		log.Fatalf("invalid tagger email, expecting: 'batman@dc.com' but received: '%s'", tag.Tagger.Email)
	}

	if tag.Fingerprint != gpgFingerprint {
		log.Fatalf("invalid fingerprint, expecting: '%s' but received: '%s'", gpgFingerprint, tag.Fingerprint)
	}

	if tag.SignedBy.Name != "batman" {
		log.Fatalf("invalid signed-by name, expecting: 'batman' but received: '%s'", tag.SignedBy.Name)
	}

	if tag.SignedBy.Email != "batman@dc.com" {
		log.Fatalf("invalid signed-by email, expecting: 'batman@dc.com' but received: '%s'", tag.SignedBy.Email)
	}
}
