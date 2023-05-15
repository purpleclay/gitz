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

	if _, err := gitc.Tag("0.1.0", git.WithSigned(), git.WithAnnotation("withsigned")); err != nil {
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

	if verif.Fingerprint != gpgFingerprint {
		log.Fatalf("invalid fingerprint, expecting: '%s' but received: '%s'", gpgFingerprint, verif.Fingerprint)
	}

	if verif.SignedBy.Name != "batman" {
		log.Fatalf("invalid signed-by name, expecting: 'batman' but received: '%s'", verif.SignedBy.Name)
	}

	if verif.SignedBy.Email != "batman@dc.com" {
		log.Fatalf("invalid signed-by email, expecting: 'batman@dc.com' but received: '%s'", verif.SignedBy.Email)
	}
}
