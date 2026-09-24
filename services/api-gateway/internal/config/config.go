package config

import "os"

type Config struct {
	HTTPAddr             string
	AuthGRPCAddr         string
	TicketGRPCAddr       string
	AssignmentGRPCAddr   string
	AuditGRPCAddr        string
	SLAGRPCAddr          string
	SearchGRPCAddr       string
	ConversationGRPCAddr string
	ProjectGRPCAddr      string
	WidgetSiteKey        string
}

func Load() Config {
	return Config{
		HTTPAddr:             getenv("GATEWAY_HTTP_ADDR", ":8080"),
		AuthGRPCAddr:         getenv("AUTH_GRPC_ADDR", "localhost:50051"),
		TicketGRPCAddr:       getenv("TICKET_GRPC_ADDR", "localhost:50052"),
		AssignmentGRPCAddr:   getenv("ASSIGNMENT_GRPC_ADDR", "localhost:50053"),
		AuditGRPCAddr:        getenv("AUDIT_GRPC_ADDR", "localhost:50054"),
		SLAGRPCAddr:          getenv("SLA_GRPC_ADDR", "localhost:50055"),
		SearchGRPCAddr:       getenv("SEARCH_GRPC_ADDR", "localhost:50057"),
		ConversationGRPCAddr: getenv("CONVERSATION_GRPC_ADDR", "localhost:50058"),
		ProjectGRPCAddr:      getenv("PROJECT_GRPC_ADDR", "localhost:50059"),
		WidgetSiteKey:        getenv("WIDGET_SITE_KEY", "demo-site"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
