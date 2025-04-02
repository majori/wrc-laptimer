package nfc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/peterhellberg/acr122u"
)

func ReadCardReader(ctx context.Context, events <-chan string) error {
	readerCtx, err := acr122u.EstablishContext()
	if err != nil {
		return err
	}

	channel := make(chan string, 1)

	log.Println("ready for smartcard events")
	go readerCtx.ServeFunc(func(c acr122u.Card) {
		hasher := sha256.New()
		hasher.Write(c.UID())

		channel <- base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	})

	return nil
}
