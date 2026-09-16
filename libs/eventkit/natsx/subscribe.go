package natsx

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/nats-io/nats.go/jetstream"
)

// Handler processes a single event. Return nil to ack.
type Handler func(ctx context.Context, ev envelope.Event) error

// SubscriberConfig configures a durable pull consumer.
type SubscriberConfig struct {
	Stream   string
	Durable  string
	Filter   string // subject filter, e.g. helpdesk.ticket.created
	Logger   *slog.Logger
}

// Subscribe starts a pull consumer loop until ctx is cancelled.
// This is a skeleton: ack on success, nak on handler error.
func Subscribe(ctx context.Context, js jetstream.JetStream, cfg SubscriberConfig, handler Handler) error {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	cons, err := js.CreateOrUpdateConsumer(ctx, cfg.Stream, jetstream.ConsumerConfig{
		Durable:       cfg.Durable,
		FilterSubject: cfg.Filter,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			msgs, err := cons.Fetch(1, jetstream.FetchMaxWait(2*time.Second))
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}
			for msg := range msgs.Messages() {
				ev, err := envelope.Unmarshal(msg.Data())
				if err != nil {
					cfg.Logger.Error("bad event payload", "error", err)
					_ = msg.Term()
					continue
				}
				if err := handler(ctx, ev); err != nil {
					cfg.Logger.Error("handler failed",
						"type", ev.Type,
						"event_id", ev.EventID,
						"error", err,
					)
					_ = msg.Nak()
					continue
				}
				_ = msg.Ack()
			}
		}
	}()

	return nil
}
