package telemetry

import (
	"context"
	"log"
	"net"
)

func StartUDPReceiver(ctx context.Context, listen string, ch chan<- any) error {
	// Validate the UDP address
	host, port, err := net.SplitHostPort(listen)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(host, port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer conn.Close()

	log.Println("listen udp:", addr)
	defer log.Println("udp closed:", addr)

	done := make(chan error, 1)
	go func() {
		// Create a buffer large enough
		b := make([]byte, 256)
		for {
			n, _, err := conn.ReadFrom(b)
			if err != nil {
				done <- err
				continue
			}

			// Process only the bytes that were read
			_, pkt, err := UnmarshalBinary(b[:n])
			if err != nil {
				done <- err
				continue
			}

			// TODO: Drop packets which are in wrong order

			ch <- pkt
		}
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
	}
	return nil
}
