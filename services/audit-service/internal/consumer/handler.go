package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Glistand/HelpDesk/libs/eventkit/envelope"
	"github.com/Glistand/HelpDesk/libs/eventkit/natsx"
	"github.com/Glistand/HelpDesk/libs/eventkit/subjects"
	"github.com/Glistand/HelpDesk/services/audit-service/internal/repository"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler struct {
	repo   *repository.Repo
	logger *slog.Logger
}

func New(repo *repository.Repo, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{repo: repo, logger: logger}
}

func (h *Handler) Handle(ctx context.Context, ev envelope.Event) error {
	title, detail, actor, ok := mapEvent(ev)
	if !ok {
		return nil
	}

	ticketID := ev.AggregateID
	if ticketID == "" {
		var p map[string]any
		_ = json.Unmarshal(ev.Payload, &p)
		if tid, ok := p["ticket_id"].(string); ok {
			ticketID = tid
		}
	}
	if ticketID == "" {
		h.logger.Warn("audit skip: no ticket_id", "event_id", ev.EventID, "type", ev.Type)
		return nil
	}

	_, err := h.repo.Insert(ctx, repository.TimelineEvent{
		ID:         uuid.NewString(),
		TicketID:   ticketID,
		EventID:    ev.EventID,
		EventType:  ev.Type,
		Title:      title,
		Detail:     detail,
		Actor:      actor,
		OccurredAt: ev.OccurredAt.UTC(),
	})
	if err != nil {
		return err
	}
	h.logger.Info("timeline appended", "ticket_id", ticketID, "type", ev.Type, "event_id", ev.EventID)
	return nil
}

func mapEvent(ev envelope.Event) (title, detail, actor string, ok bool) {
	var p map[string]any
	_ = json.Unmarshal(ev.Payload, &p)

	switch ev.Type {
	case "ticket.created":
		requester, _ := p["requester"].(string)
		return "Тикет создан", fmt.Sprintf("%v", p["title"]), requester, true
	case "ticket.assigned":
		aid, _ := p["assignee_id"].(string)
		return "Назначен агент", "assignee_id=" + aid, "assignment-service", true
	case "ticket.updated":
		status, _ := p["status"].(string)
		return "Статус изменён", "status=" + status, "", true
	case "notification.sent":
		channel, _ := p["channel"].(string)
		recipient, _ := p["recipient"].(string)
		return "Уведомление отправлено", channel + " → " + recipient, "notification-service", true
	default:
		return "", "", "", false
	}
}

func Subscribe(ctx context.Context, js jetstream.JetStream, h *Handler) error {
	return natsx.Subscribe(ctx, js, natsx.SubscriberConfig{
		Stream: subjects.StreamEvents,
		Durable: "audit-timeline",
		FilterSubjects: []string{
			subjects.TicketCreated,
			subjects.TicketAssigned,
			subjects.TicketUpdated,
			subjects.NotificationSent,
		},
		Logger: h.logger,
	}, h.Handle)
}
