package grpcserver

import (
	"context"
	"time"

	escalationv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/escalation/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/escalation-service/internal/repository"
)

type Server struct {
	escalationv1.UnimplementedEscalationServiceServer
	repo *repository.Repo
}

func New(repo *repository.Repo) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetEscalation(ctx context.Context, req *escalationv1.GetEscalationRequest) (*escalationv1.GetEscalationResponse, error) {
	if req.GetTicketId() == "" {
		return nil, statuserr.InvalidArgument("ticket_id is required")
	}
	e, err := s.repo.Get(ctx, req.GetTicketId())
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, statuserr.NotFound("escalation not found")
		}
		return nil, statuserr.FromError(err)
	}
	return &escalationv1.GetEscalationResponse{
		Escalation: &escalationv1.Escalation{
			TicketId:     e.TicketID,
			Reason:       e.Reason,
			FromAssignee: e.FromAssignee,
			ToAssignee:   e.ToAssignee,
			EscalatedAt:  e.EscalatedAt.UTC().Format(time.RFC3339),
			SourceEvent:  e.SourceEvent,
		},
	}, nil
}
