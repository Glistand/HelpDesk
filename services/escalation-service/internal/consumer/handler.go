package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	ticketv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/ticket/v1"
	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/services/escalation-service/internal/repository"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler struct {
	repo         *repository.Repo
	ticket       ticketv1.TicketServiceClient
	pub          *natsx.Publisher
	l2AssigneeID string
	logger       *slog.Logger
}

func New(
	repo *repository.Repo,
	ticket ticketv1.TicketServiceClient,
	pub *natsx.Publisher,
	l2AssigneeID string,
	logger *slog.Logger,
) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		repo:         repo,
		ticket:       ticket,
		pub:          pub,
		l2AssigneeID: l2AssigneeID,
		logger:       logger,
	}
}

func (h *Handler) Handle(ctx context.Context, ev envelope.Event) error {
	if ev.Type != "sla.breached" {
		return nil
	}

	done, err := h.repo.AlreadyProcessed(ctx, ev.EventID)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	var p map[string]any
	_ = json.Unmarshal(ev.Payload, &p)
	kind, _ := p["kind"].(string)
	// Escalate on first-response breach only (avoid double escalate on resolve breach).
	if kind != "" && kind != "fr_breach" {
		return h.repo.MarkProcessed(ctx, ev.EventID)
	}

	ticketID := ev.AggregateID
	if ticketID == "" {
		if tid, ok := p["ticket_id"].(string); ok {
			ticketID = tid
		}
	}
	if ticketID == "" {
		return nil
	}

	if _, err := h.repo.Get(ctx, ticketID); err == nil {
		return h.repo.MarkProcessed(ctx, ev.EventID)
	}

	cur, err := h.ticket.GetTicket(ctx, &ticketv1.GetTicketRequest{Id: ticketID})
	if err != nil {
		return err
	}
	from := cur.GetTicket().GetAssigneeId()
	if cur.GetTicket().GetAssigneeId() == h.l2AssigneeID {
		return h.repo.MarkProcessed(ctx, ev.EventID)
	}

	_, err = h.ticket.AssignTicket(ctx, &ticketv1.AssignTicketRequest{
		Id:         ticketID,
		AssigneeId: h.l2AssigneeID,
	})
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	inserted, err := h.repo.Save(ctx, repository.Escalation{
		TicketID:     ticketID,
		Reason:       "sla.breached:" + kind,
		FromAssignee: from,
		ToAssignee:   h.l2AssigneeID,
		EscalatedAt:  now,
		SourceEvent:  ev.EventID,
	})
	if err != nil {
		return err
	}
	if !inserted {
		return nil
	}

	out, err := envelope.New("ticket.escalated", ticketID, ev.CorrelationID, map[string]any{
		"ticket_id":     ticketID,
		"reason":        "sla.breached",
		"kind":          kind,
		"from_assignee": from,
		"to_assignee":   h.l2AssigneeID,
	})
	if err != nil {
		return err
	}
	if err := h.pub.Publish(ctx, subjects.TicketEscalated, out); err != nil {
		return err
	}

	h.logger.Info("ticket escalated",
		"ticket_id", ticketID,
		"from", from,
		"to", h.l2AssigneeID,
		"kind", kind,
	)
	return nil
}

func Subscribe(ctx context.Context, js jetstream.JetStream, h *Handler) error {
	return natsx.Subscribe(ctx, js, natsx.SubscriberConfig{
		Stream:    subjects.StreamEvents,
		Durable:   "escalation-breach",
		Filter:    subjects.SLABreached,
		EnableDLQ: true,
		Logger:    h.logger,
	}, h.Handle)
}
