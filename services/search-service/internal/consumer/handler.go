package consumer

import (
	"context"
	"log/slog"

	ticketv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/ticket/v1"
	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/services/search-service/internal/index"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler struct {
	ticket ticketv1.TicketServiceClient
	idx    *index.Store
	logger *slog.Logger
}

func New(ticket ticketv1.TicketServiceClient, idx *index.Store, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{ticket: ticket, idx: idx, logger: logger}
}

func (h *Handler) Handle(ctx context.Context, ev envelope.Event) error {
	switch ev.Type {
	case "ticket.created", "ticket.assigned", "ticket.updated", "ticket.escalated":
	default:
		return nil
	}
	ticketID := ev.AggregateID
	if ticketID == "" {
		return nil
	}

	resp, err := h.ticket.GetTicket(ctx, &ticketv1.GetTicketRequest{Id: ticketID})
	if err != nil {
		return err
	}
	t := resp.GetTicket()
	doc := index.Document{
		ID:          t.GetId(),
		Title:       t.GetTitle(),
		Description: t.GetDescription(),
		Status:      statusString(t.GetStatus()),
		Priority:    priorityString(t.GetPriority()),
		Category:    t.GetCategory(),
		Requester:   t.GetRequester(),
		AssigneeID:  t.GetAssigneeId(),
		UpdatedAt:   t.GetUpdatedAt(),
	}
	if err := h.idx.Upsert(ctx, doc); err != nil {
		return err
	}
	h.logger.Info("ticket indexed", "ticket_id", ticketID, "source", ev.Type)
	return nil
}

func statusString(s ticketv1.TicketStatus) string {
	switch s {
	case ticketv1.TicketStatus_TICKET_STATUS_NEW:
		return "new"
	case ticketv1.TicketStatus_TICKET_STATUS_OPEN:
		return "open"
	case ticketv1.TicketStatus_TICKET_STATUS_PENDING:
		return "pending"
	case ticketv1.TicketStatus_TICKET_STATUS_RESOLVED:
		return "resolved"
	case ticketv1.TicketStatus_TICKET_STATUS_CLOSED:
		return "closed"
	default:
		return "unspecified"
	}
}

func priorityString(p ticketv1.TicketPriority) string {
	switch p {
	case ticketv1.TicketPriority_TICKET_PRIORITY_LOW:
		return "low"
	case ticketv1.TicketPriority_TICKET_PRIORITY_HIGH:
		return "high"
	case ticketv1.TicketPriority_TICKET_PRIORITY_URGENT:
		return "urgent"
	default:
		return "normal"
	}
}

func Subscribe(ctx context.Context, js jetstream.JetStream, h *Handler) error {
	return natsx.Subscribe(ctx, js, natsx.SubscriberConfig{
		Stream:  subjects.StreamEvents,
		Durable: "search-indexer",
		FilterSubjects: []string{
			subjects.TicketCreated,
			subjects.TicketAssigned,
			subjects.TicketUpdated,
			subjects.TicketEscalated,
		},
		EnableDLQ: true,
		Logger:    h.logger,
	}, h.Handle)
}
