package grpcserver

import (
	"context"

	searchv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/search/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/search-service/internal/index"
)

type Server struct {
	searchv1.UnimplementedSearchServiceServer
	idx *index.Store
}

func New(idx *index.Store) *Server {
	return &Server{idx: idx}
}

func (s *Server) SearchTickets(ctx context.Context, req *searchv1.SearchTicketsRequest) (*searchv1.SearchTicketsResponse, error) {
	limit := int64(req.GetLimit())
	hits, total, err := s.idx.Search(ctx, req.GetQuery(), limit)
	if err != nil {
		return nil, statuserr.FromError(err)
	}
	out := make([]*searchv1.SearchHit, 0, len(hits))
	for _, h := range hits {
		out = append(out, &searchv1.SearchHit{
			Id:          h.ID,
			Title:       h.Title,
			Description: h.Description,
			Status:      h.Status,
			Priority:    h.Priority,
			Category:    h.Category,
			Requester:   h.Requester,
			AssigneeId:  h.AssigneeID,
			UpdatedAt:   h.UpdatedAt,
		})
	}
	return &searchv1.SearchTicketsResponse{Hits: out, EstimatedTotal: total}, nil
}
