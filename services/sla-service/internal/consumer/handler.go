package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/policy"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/repository"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/timers"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler struct {
	repo   *repository.Repo
	timers *timers.Store
	pub    *natsx.Publisher
	pol    policy.Policy
	logger *slog.Logger
}

func New(repo *repository.Repo, t *timers.Store, pub *natsx.Publisher, pol policy.Policy, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{repo: repo, timers: t, pub: pub, pol: pol, logger: logger}
}

func (h *Handler) Handle(ctx context.Context, ev envelope.Event) error {
	switch ev.Type {
	case "ticket.created", "ticket.assigned":
		return h.arm(ctx, ev)
	case "ticket.updated":
		return h.onUpdated(ctx, ev)
	default:
		return nil
	}
}

func (h *Handler) arm(ctx context.Context, ev envelope.Event) error {
	done, err := h.repo.AlreadyProcessed(ctx, ev.EventID)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	ticketID := ev.AggregateID
	if ticketID == "" {
		return nil
	}

	// Only arm once per ticket (first created/assigned wins).
	if _, err := h.repo.Get(ctx, ticketID); err == nil {
		return h.repo.MarkProcessed(ctx, ev.EventID)
	}

	now := time.Now().UTC()
	frDue := now.Add(h.pol.FirstResponse)
	resolveDue := now.Add(h.pol.Resolve)
	frWarn := h.pol.WarnAt(now, frDue)
	resolveWarn := h.pol.WarnAt(now, resolveDue)

	sla := repository.TicketSLA{
		TicketID:         ticketID,
		State:            "ok",
		Policy:           h.pol.Name,
		FirstResponseDue: frDue,
		ResolveDue:       resolveDue,
		CorrelationID:    ev.CorrelationID,
	}
	if err := h.repo.Upsert(ctx, sla); err != nil {
		return err
	}
	_ = h.repo.ClearFired(ctx, ticketID)

	if err := h.timers.Schedule(ctx, ticketID, map[string]time.Time{
		timers.KindFRWarn:        frWarn,
		timers.KindFRBreach:      frDue,
		timers.KindResolveWarn:   resolveWarn,
		timers.KindResolveBreach: resolveDue,
	}); err != nil {
		return err
	}

	h.logger.Info("sla armed",
		"ticket_id", ticketID,
		"first_response_due", frDue.Format(time.RFC3339),
		"resolve_due", resolveDue.Format(time.RFC3339),
		"policy", h.pol.Name,
	)
	return h.repo.MarkProcessed(ctx, ev.EventID)
}

func (h *Handler) onUpdated(ctx context.Context, ev envelope.Event) error {
	done, err := h.repo.AlreadyProcessed(ctx, ev.EventID)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	var p map[string]any
	_ = json.Unmarshal(ev.Payload, &p)
	status, _ := p["status"].(string)
	if status != "resolved" && status != "closed" {
		return h.repo.MarkProcessed(ctx, ev.EventID)
	}

	ticketID := ev.AggregateID
	now := time.Now().UTC()
	_ = h.timers.CancelTicket(ctx, ticketID)
	_ = h.repo.Cancel(ctx, ticketID, now)
	h.logger.Info("sla cancelled", "ticket_id", ticketID, "status", status)
	return h.repo.MarkProcessed(ctx, ev.EventID)
}

func Subscribe(ctx context.Context, js jetstream.JetStream, h *Handler) error {
	return natsx.Subscribe(ctx, js, natsx.SubscriberConfig{
		Stream:  subjects.StreamEvents,
		Durable: "sla-lifecycle",
		FilterSubjects: []string{
			subjects.TicketCreated,
			subjects.TicketAssigned,
			subjects.TicketUpdated,
		},
		EnableDLQ: true,
		Logger:    h.logger,
	}, h.Handle)
}
