package grpcserver

import (
	"context"
	"time"

	assignmentv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/assignment/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/assignment-service/internal/repository"
)

type Server struct {
	assignmentv1.UnimplementedAssignmentServiceServer
	repo *repository.Repo
}

func New(repo *repository.Repo) *Server {
	return &Server{repo: repo}
}

func (s *Server) GetAssignment(ctx context.Context, req *assignmentv1.GetAssignmentRequest) (*assignmentv1.GetAssignmentResponse, error) {
	a, err := s.repo.GetByTicket(ctx, req.GetTicketId())
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, statuserr.NotFound("assignment not found")
		}
		return nil, statuserr.FromError(err)
	}
	return &assignmentv1.GetAssignmentResponse{
		Assignment: &assignmentv1.Assignment{
			TicketId:     a.TicketID,
			AssigneeId:   a.AssigneeID,
			AssigneeName: a.AssigneeName,
			AssignedAt:   a.AssignedAt.UTC().Format(time.RFC3339),
		},
	}, nil
}
