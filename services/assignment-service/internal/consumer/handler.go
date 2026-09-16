package consumer

import (
	"context"
	"log/slog"
	"time"

	ticketv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/ticket/v1"
	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/services/assignment-service/internal/repository"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler struct {
	repo   *repository.Repo
	ticket ticketv1.TicketServiceClient
	logger *slog.Logger
}

func New(repo *repository.Repo, ticket ticketv1.TicketServiceClient, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{repo: repo, ticket: ticket, logger: logger}
}

func (h *Handler) Handle(ctx context.Context, ev envelope.Event) error {
	if ev.Type != "ticket.created" {
		return nil
	}

	done, err := h.repo.AlreadyProcessed(ctx, ev.EventID)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	ticketID := ev.AggregateID
	if ticketID == "" {
		h.logger.Warn("missing aggregate_id", "event_id", ev.EventID)
		return nil
	}

	if _, err := h.repo.GetByTicket(ctx, ticketID); err == nil {
		return h.repo.MarkProcessed(ctx, ev.EventID)
	}

	agent, err := h.repo.NextL1Agent(ctx)
	if err != nil {
		return err
	}

	_, err = h.ticket.AssignTicket(ctx, &ticketv1.AssignTicketRequest{
		Id:         ticketID,
		AssigneeId: agent.ID,
	})
	if err != nil {
		return err
	}

	_, err = h.repo.SaveAssignment(ctx, repository.Assignment{
		TicketID:     ticketID,
		AssigneeID:   agent.ID,
		AssigneeName: agent.Name,
		AssignedAt:   time.Now().UTC(),
		SourceEvent:  ev.EventID,
	})
	if err != nil {
		return err
	}

	h.logger.Info("ticket assigned",
		"ticket_id", ticketID,
		"assignee_id", agent.ID,
		"assignee_name", agent.Name,
		"event_id", ev.EventID,
	)
	return nil
}

func Subscribe(ctx context.Context, js jetstream.JetStream, h *Handler) error {
	return natsx.Subscribe(ctx, js, natsx.SubscriberConfig{
		Stream:    subjects.StreamEvents,
		Durable:   "assignment-created",
		Filter:    subjects.TicketCreated,
		EnableDLQ: true,
		Logger:    h.logger,
	}, h.Handle)
}
