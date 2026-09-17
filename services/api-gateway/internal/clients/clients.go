package clients

import (
	"context"

	assignmentv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/assignment/v1"
	auditv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/audit/v1"
	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
	searchv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/search/v1"
	slav1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/sla/v1"
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
	SLA        slav1.SLAServiceClient
	Search     searchv1.SearchServiceClient
	conns      []*grpc.ClientConn
}

func Dial(ctx context.Context, authAddr, ticketAddr, assignmentAddr, auditAddr, slaAddr, searchAddr string) (*Clients, error) {
	dial := func(addr string) (*grpc.ClientConn, error) {
		opts := append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, grpckit.DefaultClientOptions()...)
		return grpc.NewClient(addr, opts...)
	}

	addrs := []string{authAddr, ticketAddr, assignmentAddr, auditAddr, slaAddr, searchAddr}
	conns := make([]*grpc.ClientConn, 0, len(addrs))
	for _, addr := range addrs {
		c, err := dial(addr)
		if err != nil {
			for _, x := range conns {
				_ = x.Close()
			}
			return nil, err
		}
		conns = append(conns, c)
	}

	return &Clients{
		Auth:       authv1.NewAuthServiceClient(conns[0]),
		Ticket:     ticketv1.NewTicketServiceClient(conns[1]),
		Assignment: assignmentv1.NewAssignmentServiceClient(conns[2]),
		Audit:      auditv1.NewAuditServiceClient(conns[3]),
		SLA:        slav1.NewSLAServiceClient(conns[4]),
		Search:     searchv1.NewSearchServiceClient(conns[5]),
		conns:      conns,
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.conns {
		_ = conn.Close()
	}
}
