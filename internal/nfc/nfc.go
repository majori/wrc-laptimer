package nfc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/peterhellberg/acr122u"
)

func ListenForCardEvents(ctx context.Context, events chan<- string) error {
	readerCtx, err := acr122u.EstablishContext()
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer readerCtx.Release()

	log.Println("ready for smartcard events")
	return readerCtx.ServeFunc(func(c acr122u.Card) {
		hasher := sha256.New()
		hasher.Write(c.UID())

		events <- hex.EncodeToString(hasher.Sum(nil))
	})
}
