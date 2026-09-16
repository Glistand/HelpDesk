package natsx

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const defaultURL = "nats://localhost:4222"

// Connect opens a NATS connection and JetStream context.
// URL is taken from NATS_URL env or defaults to nats://localhost:4222.
func Connect(ctx context.Context) (*nats.Conn, jetstream.JetStream, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = defaultURL
	}

	nc, err := nats.Connect(url,
		nats.Name("helpdesk"),
		nats.Timeout(5*time.Second),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, nil, fmt.Errorf("jetstream: %w", err)
	}

	if err := waitConnected(ctx, nc); err != nil {
		nc.Close()
		return nil, nil, err
	}

	return nc, js, nil
}

func waitConnected(ctx context.Context, nc *nats.Conn) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(5 * time.Second)
	}
	for time.Now().Before(deadline) {
		if nc.IsConnected() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
	if nc.IsConnected() {
		return nil
	}
	return fmt.Errorf("nats not connected to %s", nc.ConnectedUrl())
}
