package clients

import (
	"context"

	assignmentv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/assignment/v1"
	auditv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/audit/v1"
	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
	ticketv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/ticket/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Auth       authv1.AuthServiceClient
	Ticket     ticketv1.TicketServiceClient
	Assignment assignmentv1.AssignmentServiceClient
	Audit      auditv1.AuditServiceClient
	conns      []*grpc.ClientConn
}

func Dial(ctx context.Context, authAddr, ticketAddr, assignmentAddr, auditAddr string) (*Clients, error) {
	dial := func(addr string) (*grpc.ClientConn, error) {
		return grpc.NewClient(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpckit.DefaultUnaryClientInterceptors(),
		)
	}

	authConn, err := dial(authAddr)
	if err != nil {
		return nil, err
	}
	ticketConn, err := dial(ticketAddr)
	if err != nil {
		_ = authConn.Close()
		return nil, err
	}
	assignmentConn, err := dial(assignmentAddr)
	if err != nil {
		_ = authConn.Close()
		_ = ticketConn.Close()
		return nil, err
	}
	auditConn, err := dial(auditAddr)
	if err != nil {
		_ = authConn.Close()
		_ = ticketConn.Close()
		_ = assignmentConn.Close()
		return nil, err
	}

	return &Clients{
		Auth:       authv1.NewAuthServiceClient(authConn),
		Ticket:     ticketv1.NewTicketServiceClient(ticketConn),
		Assignment: assignmentv1.NewAssignmentServiceClient(assignmentConn),
		Audit:      auditv1.NewAuditServiceClient(auditConn),
		conns:      []*grpc.ClientConn{authConn, ticketConn, assignmentConn, auditConn},
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.conns {
		_ = conn.Close()
	}
}
