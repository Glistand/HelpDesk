package clients

import (
	"context"

	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
	ticketv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/ticket/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Auth   authv1.AuthServiceClient
	Ticket ticketv1.TicketServiceClient
	conns  []*grpc.ClientConn
}

func Dial(ctx context.Context, authAddr, ticketAddr string) (*Clients, error) {
	authConn, err := grpc.NewClient(authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpckit.DefaultUnaryClientInterceptors(),
	)
	if err != nil {
		return nil, err
	}
	ticketConn, err := grpc.NewClient(ticketAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpckit.DefaultUnaryClientInterceptors(),
	)
	if err != nil {
		_ = authConn.Close()
		return nil, err
	}
	return &Clients{
		Auth:   authv1.NewAuthServiceClient(authConn),
		Ticket: ticketv1.NewTicketServiceClient(ticketConn),
		conns:  []*grpc.ClientConn{authConn, ticketConn},
	}, nil
}

func (c *Clients) Close() {
	for _, conn := range c.conns {
		_ = conn.Close()
	}
}
