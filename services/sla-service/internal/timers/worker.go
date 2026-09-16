package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/repository"
)

// Worker polls Redis for due SLA timers and publishes warn/breach events.
type Worker struct {
	store    *Store
	repo     *repository.Repo
	pub      *natsx.Publisher
	interval time.Duration
	logger   *slog.Logger
}

func NewWorker(store *Store, repo *repository.Repo, pub *natsx.Publisher, interval time.Duration, logger *slog.Logger) *Worker {
	if logger == nil {
		logger = slog.Default()
	}
	return &Worker{store: store, repo: repo, pub: pub, interval: interval, logger: logger}
}

func (w *Worker) Run(ctx context.Context) {
	t := time.NewTicker(w.interval)
	defer t.Stop()
	for {
		w.tick(ctx)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	due, err := w.store.PopDue(ctx, time.Now().UTC(), 50)
	if err != nil {
		w.logger.Error("sla pop due failed", "error", err)
		return
	}
	for _, d := range due {
		if err := w.fire(ctx, d); err != nil {
			w.logger.Error("sla fire failed",
				"ticket_id", d.TicketID,
				"kind", d.Kind,
				"error", err,
			)
		}
	}
}

func (w *Worker) fire(ctx context.Context, d Due) error {
	sla, err := w.repo.Get(ctx, d.TicketID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil
		}
		return err
	}
	if sla.CancelledAt != nil {
		return nil
	}

	first, err := w.repo.MarkFired(ctx, d.TicketID, d.Kind)
	if err != nil {
		return err
	}
	if !first {
		return nil
	}

	now := time.Now().UTC()
	switch d.Kind {
	case KindFRWarn, KindResolveWarn:
		if err := w.repo.MarkWarned(ctx, d.TicketID, now); err != nil {
			return err
		}
		ev, err := envelope.New("sla.warned", d.TicketID, sla.CorrelationID, map[string]any{
			"ticket_id": d.TicketID,
			"kind":      d.Kind,
			"due_at":    d.DueAt.Format(time.RFC3339),
			"policy":    sla.Policy,
		})
		if err != nil {
			return err
		}
		w.logger.Info("sla warned", "ticket_id", d.TicketID, "kind", d.Kind)
		return w.pub.Publish(ctx, subjects.SLAWarned, ev)

	case KindFRBreach, KindResolveBreach:
		if err := w.repo.MarkBreached(ctx, d.TicketID, now); err != nil {
			return err
		}
		ev, err := envelope.New("sla.breached", d.TicketID, sla.CorrelationID, map[string]any{
			"ticket_id": d.TicketID,
			"kind":      d.Kind,
			"due_at":    d.DueAt.Format(time.RFC3339),
			"policy":    sla.Policy,
		})
		if err != nil {
			return err
		}
		w.logger.Info("sla breached", "ticket_id", d.TicketID, "kind", d.Kind)
		return w.pub.Publish(ctx, subjects.SLABreached, ev)
	}
	return nil
}
