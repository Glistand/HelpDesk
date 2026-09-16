package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync/atomic"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/services/notification-service/internal/repository"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler struct {
	repo       *repository.Repo
	pub        *natsx.Publisher
	logger     *slog.Logger
	failFirstN int
	attempts   atomic.Int64
}

func New(repo *repository.Repo, pub *natsx.Publisher, failFirstN int, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{repo: repo, pub: pub, failFirstN: failFirstN, logger: logger}
}

func (h *Handler) Handle(ctx context.Context, ev envelope.Event) error {
	switch ev.Type {
	case "ticket.created", "ticket.assigned":
	default:
		return nil
	}

	done, err := h.repo.AlreadyProcessed(ctx, ev.EventID)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	n := h.attempts.Add(1)
	if h.failFirstN > 0 && int(n) <= h.failFirstN {
		return fmt.Errorf("simulated notify failure attempt=%d", n)
	}

	ticketID := ev.AggregateID
	channel := "email"
	recipient := "agent@helpdesk.local"
	subj := fmt.Sprintf("Ticket %s: %s", ticketID, ev.Type)
	body := fmt.Sprintf("Mock notification for %s (correlation=%s)", ticketID, ev.CorrelationID)

	var payload map[string]any
	_ = json.Unmarshal(ev.Payload, &payload)
	if aid, ok := payload["assignee_id"].(string); ok && aid != "" {
		recipient = aid + "@helpdesk.local"
	}

	h.logger.Info("mock notify",
		"channel", channel,
		"recipient", recipient,
		"ticket_id", ticketID,
		"source_type", ev.Type,
	)

	sent, err := envelope.New("notification.sent", ticketID, ev.CorrelationID, map[string]any{
		"ticket_id":    ticketID,
		"channel":      channel,
		"recipient":    recipient,
		"subject":      subj,
		"source_event": ev.EventID,
		"source_type":  ev.Type,
	})
	if err != nil {
		return err
	}

	if err := h.pub.Publish(ctx, subjects.NotificationSent, sent); err != nil {
		return err
	}

	return h.repo.Save(ctx, repository.Notification{
		EventID:   ev.EventID,
		TicketID:  ticketID,
		Channel:   channel,
		Recipient: recipient,
		Subject:   subj,
		Body:      body,
		Status:    "sent",
		Attempts:  int(n),
	})
}

func Subscribe(ctx context.Context, js jetstream.JetStream, h *Handler) error {
	return natsx.Subscribe(ctx, js, natsx.SubscriberConfig{
		Stream:         subjects.StreamEvents,
		Durable:        "notification-events",
		FilterSubjects: []string{subjects.TicketCreated, subjects.TicketAssigned},
		Logger:         h.logger,
	}, h.Handle)
}
