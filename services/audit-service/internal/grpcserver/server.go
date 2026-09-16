package grpcserver

import (
	"context"
	"time"

	auditv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/audit/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/audit-service/internal/repository"
)

type Server struct {
	auditv1.UnimplementedAuditServiceServer
	repo *repository.Repo
}

func New(repo *repository.Repo) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetTimeline(ctx context.Context, req *auditv1.GetTimelineRequest) (*auditv1.GetTimelineResponse, error) {
	if req.GetTicketId() == "" {
		return nil, statuserr.InvalidArgument("ticket_id is required")
	}
	events, err := s.repo.ListByTicket(ctx, req.GetTicketId())
	if err != nil {
		return nil, statuserr.FromError(err)
	}
	out := make([]*auditv1.TimelineEvent, 0, len(events))
	for _, e := range events {
		out = append(out, &auditv1.TimelineEvent{
			Id:         e.ID,
			TicketId:   e.TicketID,
			EventId:    e.EventID,
			EventType:  e.EventType,
			Title:      e.Title,
			Detail:     e.Detail,
			Actor:      e.Actor,
			OccurredAt: e.OccurredAt.UTC().Format(time.RFC3339),
		})
	}
	return &auditv1.GetTimelineResponse{Events: out}, nil
}
