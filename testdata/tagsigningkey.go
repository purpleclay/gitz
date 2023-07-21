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

	if _, err := gitc.Tag("0.1.0",
		git.WithSigningKey(gpgPublicKeyID),
		git.WithAnnotation("withsigningkey")); err != nil {
		log.Fatal(err.Error())
	}

	verif, err := gitc.VerifyTag("0.1.0")
	if err != nil {
		log.Fatal(err.Error())
	}

	if verif.Tagger.Name != "batman" {
		log.Fatalf("invalid tagger name, expecting: 'batman' but received: '%s'", verif.Tagger.Name)
	}

	if verif.Tagger.Email != "batman@dc.com" {
		log.Fatalf("invalid tagger email, expecting: 'batman@dc.com' but received: '%s'", verif.Tagger.Email)
	}

	if verif.Annotation != "withsigningkey" {
		log.Fatalf("invalid annotation, expecting: 'withsigningkey' but received: '%s'", verif.Annotation)
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
