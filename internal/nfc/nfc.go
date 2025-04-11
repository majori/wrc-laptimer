package nfc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"

	"github.com/peterhellberg/acr122u"
)

func ListenForCardEvents(ctx context.Context, events chan<- string) error {
	readerCtx, err := acr122u.EstablishContext()
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer readerCtx.Release()

	slog.Info("listening smartcard events")
	return readerCtx.ServeFunc(func(c acr122u.Card) {
		// Hash the UID of the card before sending it to the channel
		hasher := sha256.New()
		hasher.Write(c.UID())

		events <- hex.EncodeToString(hasher.Sum(nil))
	})
}
