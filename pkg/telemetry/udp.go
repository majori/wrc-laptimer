package telemetry

import (
	"context"
	"log/slog"
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

	slog.Info("listening udp for telemetry", "address", addr)

	done := make(chan error, 1)
	go func() {
		var latestPacketID uint64
		b := make([]byte, 256)

		for {
			n, _, err := conn.ReadFrom(b)
			if err != nil {
				done <- err
				continue
			}

			// Process only the bytes that were read
			header, pkt, err := UnmarshalBinary(b[:n])
			if err != nil {
				done <- err
				continue
			}

			// Ignore packets which come in wrong order
			if header.PacketUid < latestPacketID && header.PacketUid != 0 {
				continue
			}

			latestPacketID = header.PacketUid
			ch <- pkt
		}
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return nil
	}
}
