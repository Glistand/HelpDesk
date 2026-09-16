package outbox

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/services/ticket-service/internal/repository"
)

type Publisher struct {
	repo     *repository.TicketRepo
	nats     *natsx.Publisher
	interval time.Duration
	logger   *slog.Logger
}

func NewPublisher(repo *repository.TicketRepo, nats *natsx.Publisher, interval time.Duration, logger *slog.Logger) *Publisher {
	if logger == nil {
		logger = slog.Default()
	}
	return &Publisher{repo: repo, nats: nats, interval: interval, logger: logger}
}

func (p *Publisher) Run(ctx context.Context) {
	t := time.NewTicker(p.interval)
	defer t.Stop()
	for {
		p.flush(ctx)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

func (p *Publisher) flush(ctx context.Context) {
	events, err := p.repo.ClaimUnpublished(ctx, 50)
	if err != nil {
		p.logger.Error("outbox claim failed", "error", err)
		return
	}
	for _, e := range events {
		ev := envelope.Event{
			EventID:       e.EventID,
			CorrelationID: e.CorrelationID,
			OccurredAt:    e.CreatedAt.UTC(),
			Type:          e.EventType,
			AggregateID:   e.AggregateID,
			Payload:       json.RawMessage(e.Payload),
		}
		if err := p.nats.Publish(ctx, e.Subject, ev); err != nil {
			p.logger.Error("outbox publish failed",
				"event_id", e.EventID,
				"subject", e.Subject,
				"error", err,
			)
			continue
		}
		if err := p.repo.MarkPublished(ctx, e.ID); err != nil {
			p.logger.Error("outbox mark published failed", "id", e.ID, "error", err)
		}
	}
}
