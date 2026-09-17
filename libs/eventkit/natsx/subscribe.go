package natsx

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/nats-io/nats.go/jetstream"
)

// Handler processes a single event. Return nil to ack.
type Handler func(ctx context.Context, ev envelope.Event) error

// SubscriberConfig configures a durable pull consumer.
type SubscriberConfig struct {
	Stream         string
	Durable        string
	Filter         string   // single subject filter
	FilterSubjects []string // multi-subject filter (preferred when set)
	MaxDeliver     int      // JetStream max delivery attempts before DLQ (default 5)
	EnableDLQ      bool     // publish to HELP_DESK_DLQ after MaxDeliver
	Logger         *slog.Logger
}

// Subscribe starts a pull consumer loop until ctx is cancelled.
func Subscribe(ctx context.Context, js jetstream.JetStream, cfg SubscriberConfig, handler Handler) error {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	maxDeliver := cfg.MaxDeliver
	if maxDeliver <= 0 {
		maxDeliver = 5
	}

	consCfg := jetstream.ConsumerConfig{
		Durable:       cfg.Durable,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    maxDeliver,
		AckWait:       30 * time.Second,
	}
	if len(cfg.FilterSubjects) > 0 {
		consCfg.FilterSubjects = cfg.FilterSubjects
	} else if cfg.Filter != "" {
		consCfg.FilterSubject = cfg.Filter
	}

	cons, err := js.CreateOrUpdateConsumer(ctx, cfg.Stream, consCfg)
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

				meta, _ := msg.Metadata()
				delivered := uint64(1)
				if meta != nil {
					delivered = meta.NumDelivered
				}

				hctx, span := startConsumeSpan(ctx, ev, msg.Subject())
				err = handler(hctx, ev)
				endConsumeSpan(span, err)
				if err != nil {
					cfg.Logger.Error("handler failed",
						"type", ev.Type,
						"event_id", ev.EventID,
						"correlation_id", ev.CorrelationID,
						"delivered", delivered,
						"max_deliver", maxDeliver,
						"error", err,
					)
					if int(delivered) >= maxDeliver {
						if cfg.EnableDLQ {
							if dlqErr := publishDLQ(ctx, js, cfg.Durable, msg.Subject(), ev, err); dlqErr != nil {
								cfg.Logger.Error("dlq publish failed", "error", dlqErr)
							} else {
								cfg.Logger.Warn("moved to DLQ",
									"durable", cfg.Durable,
									"event_id", ev.EventID,
									"subject", msg.Subject(),
								)
							}
						}
						_ = msg.Term()
						continue
					}
					_ = msg.Nak()
					continue
				}
				_ = msg.Ack()
			}
		}
	}()

	return nil
}

func publishDLQ(ctx context.Context, js jetstream.JetStream, durable, originalSubject string, ev envelope.Event, cause error) error {
	payload := map[string]any{
		"original_subject": originalSubject,
		"durable":          durable,
		"error":            cause.Error(),
		"event":            ev,
	}
	dlqEv, err := envelope.New("dlq.poison", ev.AggregateID, ev.CorrelationID, payload)
	if err != nil {
		return err
	}
	// Preserve original event_id linkage in aggregate; use fresh dlq event id.
	subject := subjects.DLQSubject(originalSubject)
	data, err := dlqEv.Marshal()
	if err != nil {
		return err
	}
	_, err = js.Publish(ctx, subject, data, jetstream.WithMsgID(dlqEv.EventID))
	return err
}
