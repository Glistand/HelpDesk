package natsx

import (
	"context"
	"fmt"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/nats-io/nats.go/jetstream"
)

// Publisher publishes Event envelopes to JetStream subjects.
type Publisher struct {
	js jetstream.JetStream
}

// NewPublisher wraps a JetStream handle.
func NewPublisher(js jetstream.JetStream) *Publisher {
	return &Publisher{js: js}
}

// Publish marshals and publishes an event.
// Msg-Id is set to event_id for JetStream de-duplication.
func (p *Publisher) Publish(ctx context.Context, subject string, ev envelope.Event) error {
	data, err := ev.Marshal()
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	_, err = p.js.Publish(ctx, subject, data, jetstream.WithMsgID(ev.EventID))
	if err != nil {
		return fmt.Errorf("publish %s: %w", subject, err)
	}
	return nil
}
