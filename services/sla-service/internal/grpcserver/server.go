package grpcserver

import (
	"context"
	"time"

	slav1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/sla/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/sla-service/internal/repository"
)

type Server struct {
	slav1.UnimplementedSLAServiceServer
	repo *repository.Repo
}

func New(repo *repository.Repo) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetSLA(ctx context.Context, req *slav1.GetSLARequest) (*slav1.GetSLAResponse, error) {
	if req.GetTicketId() == "" {
		return nil, statuserr.InvalidArgument("ticket_id is required")
	}
	sla, err := s.repo.Get(ctx, req.GetTicketId())
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, statuserr.NotFound("sla not found")
		}
		return nil, statuserr.FromError(err)
	}
	out := &slav1.SLA{
		TicketId:         sla.TicketID,
		State:            toProtoState(sla.State),
		FirstResponseDue: sla.FirstResponseDue.UTC().Format(time.RFC3339),
		ResolveDue:       sla.ResolveDue.UTC().Format(time.RFC3339),
		Policy:           sla.Policy,
	}
	if sla.WarnedAt != nil {
		out.WarnedAt = sla.WarnedAt.UTC().Format(time.RFC3339)
	}
	if sla.BreachedAt != nil {
		out.BreachedAt = sla.BreachedAt.UTC().Format(time.RFC3339)
	}
	return &slav1.GetSLAResponse{Sla: out}, nil
}

func toProtoState(s string) slav1.SLAState {
	switch s {
	case "ok":
		return slav1.SLAState_SLA_STATE_OK
	case "warning":
		return slav1.SLAState_SLA_STATE_WARNING
	case "breached":
		return slav1.SLAState_SLA_STATE_BREACHED
	case "cancelled":
		return slav1.SLAState_SLA_STATE_CANCELLED
	default:
		return slav1.SLAState_SLA_STATE_UNSPECIFIED
	}
}
