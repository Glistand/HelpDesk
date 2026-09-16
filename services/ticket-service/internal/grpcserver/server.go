package grpcserver

import (
	"context"
	"strings"
	"time"

	ticketv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/ticket/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit/metadata"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/ticket-service/internal/domain"
	"github.com/Glistand/HelpDesk/services/ticket-service/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	ticketv1.UnimplementedTicketServiceServer
	repo *repository.TicketRepo
}

func New(repo *repository.TicketRepo) *Server {
	return &Server{repo: repo}
}

func (s *Server) CreateTicket(ctx context.Context, req *ticketv1.CreateTicketRequest) (*ticketv1.CreateTicketResponse, error) {
	t, err := s.repo.Create(ctx, repository.CreateInput{
		Title:         req.GetTitle(),
		Description:   req.GetDescription(),
		Priority:      fromProtoPriority(req.GetPriority()),
		Category:      req.GetCategory(),
		Requester:     req.GetRequester(),
		CorrelationID: metadata.CorrelationIDFromContext(ctx),
	})
	if err != nil {
		if err == repository.ErrEmptyTitle {
			return nil, statuserr.InvalidArgument("title is required")
		}
		return nil, statuserr.FromError(err)
	}
	return &ticketv1.CreateTicketResponse{Ticket: toProto(t)}, nil
}

func (s *Server) GetTicket(ctx context.Context, req *ticketv1.GetTicketRequest) (*ticketv1.GetTicketResponse, error) {
	t, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, statuserr.NotFound("ticket not found")
		}
		return nil, statuserr.FromError(err)
	}
	return &ticketv1.GetTicketResponse{Ticket: toProto(t)}, nil
}

func (s *Server) ListTickets(ctx context.Context, req *ticketv1.ListTicketsRequest) (*ticketv1.ListTicketsResponse, error) {
	var st domain.Status
	if req.GetStatus() != ticketv1.TicketStatus_TICKET_STATUS_UNSPECIFIED {
		st = fromProtoStatus(req.GetStatus())
	}
	limit := int(req.GetPageSize())
	list, err := s.repo.List(ctx, st, req.GetAssigneeId(), limit)
	if err != nil {
		return nil, statuserr.FromError(err)
	}
	out := make([]*ticketv1.Ticket, 0, len(list))
	for _, t := range list {
		out = append(out, toProto(t))
	}
	return &ticketv1.ListTicketsResponse{Tickets: out}, nil
}

func (s *Server) UpdateTicketStatus(ctx context.Context, req *ticketv1.UpdateTicketStatusRequest) (*ticketv1.UpdateTicketStatusResponse, error) {
	if req.GetStatus() == ticketv1.TicketStatus_TICKET_STATUS_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}
	t, err := s.repo.UpdateStatus(ctx, req.GetId(), fromProtoStatus(req.GetStatus()), metadata.CorrelationIDFromContext(ctx))
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, statuserr.NotFound("ticket not found")
		}
		return nil, statuserr.FromError(err)
	}
	return &ticketv1.UpdateTicketStatusResponse{Ticket: toProto(t)}, nil
}

func (s *Server) AssignTicket(ctx context.Context, req *ticketv1.AssignTicketRequest) (*ticketv1.AssignTicketResponse, error) {
	if strings.TrimSpace(req.GetAssigneeId()) == "" {
		return nil, status.Error(codes.InvalidArgument, "assignee_id is required")
	}
	t, err := s.repo.Assign(ctx, req.GetId(), req.GetAssigneeId(), metadata.CorrelationIDFromContext(ctx))
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, statuserr.NotFound("ticket not found")
		}
		if err == repository.ErrEmptyAssignee {
			return nil, status.Error(codes.InvalidArgument, "assignee_id is required")
		}
		return nil, statuserr.FromError(err)
	}
	return &ticketv1.AssignTicketResponse{Ticket: toProto(t)}, nil
}

func toProto(t domain.Ticket) *ticketv1.Ticket {
	return &ticketv1.Ticket{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      toProtoStatus(t.Status),
		Priority:    toProtoPriority(t.Priority),
		Category:    t.Category,
		Requester:   t.Requester,
		AssigneeId:  t.AssigneeID,
		CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toProtoStatus(s domain.Status) ticketv1.TicketStatus {
	switch s {
	case domain.StatusNew:
		return ticketv1.TicketStatus_TICKET_STATUS_NEW
	case domain.StatusOpen:
		return ticketv1.TicketStatus_TICKET_STATUS_OPEN
	case domain.StatusPending:
		return ticketv1.TicketStatus_TICKET_STATUS_PENDING
	case domain.StatusResolved:
		return ticketv1.TicketStatus_TICKET_STATUS_RESOLVED
	case domain.StatusClosed:
		return ticketv1.TicketStatus_TICKET_STATUS_CLOSED
	default:
		return ticketv1.TicketStatus_TICKET_STATUS_UNSPECIFIED
	}
}

func fromProtoStatus(s ticketv1.TicketStatus) domain.Status {
	switch s {
	case ticketv1.TicketStatus_TICKET_STATUS_NEW:
		return domain.StatusNew
	case ticketv1.TicketStatus_TICKET_STATUS_OPEN:
		return domain.StatusOpen
	case ticketv1.TicketStatus_TICKET_STATUS_PENDING:
		return domain.StatusPending
	case ticketv1.TicketStatus_TICKET_STATUS_RESOLVED:
		return domain.StatusResolved
	case ticketv1.TicketStatus_TICKET_STATUS_CLOSED:
		return domain.StatusClosed
	default:
		return ""
	}
}

func toProtoPriority(p domain.Priority) ticketv1.TicketPriority {
	switch p {
	case domain.PriorityLow:
		return ticketv1.TicketPriority_TICKET_PRIORITY_LOW
	case domain.PriorityNormal:
		return ticketv1.TicketPriority_TICKET_PRIORITY_NORMAL
	case domain.PriorityHigh:
		return ticketv1.TicketPriority_TICKET_PRIORITY_HIGH
	case domain.PriorityUrgent:
		return ticketv1.TicketPriority_TICKET_PRIORITY_URGENT
	default:
		return ticketv1.TicketPriority_TICKET_PRIORITY_UNSPECIFIED
	}
}

func fromProtoPriority(p ticketv1.TicketPriority) domain.Priority {
	switch p {
	case ticketv1.TicketPriority_TICKET_PRIORITY_LOW:
		return domain.PriorityLow
	case ticketv1.TicketPriority_TICKET_PRIORITY_HIGH:
		return domain.PriorityHigh
	case ticketv1.TicketPriority_TICKET_PRIORITY_URGENT:
		return domain.PriorityUrgent
	case ticketv1.TicketPriority_TICKET_PRIORITY_NORMAL:
		return domain.PriorityNormal
	default:
		return domain.PriorityNormal
	}
}
