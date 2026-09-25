package config

import "os"

type Config struct {
	GRPCAddr       string
	DatabaseURL    string
	NATSURL        string
	TicketGRPCAddr string
	L2AssigneeID   string
}

func Load() Config {
	return Config{
		GRPCAddr:       getenv("ESCALATION_GRPC_ADDR", ":50056"),
		DatabaseURL:    getenv("ESCALATION_DATABASE_URL", "postgres://helpdesk:helpdesk@localhost:5432/escalation?sslmode=disable"),
		NATSURL:        getenv("NATS_URL", "nats://localhost:4222"),
		TicketGRPCAddr: getenv("TICKET_GRPC_ADDR", "localhost:50052"),
		L2AssigneeID:   getenv("ESCALATION_L2_ASSIGNEE", "a-1"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
