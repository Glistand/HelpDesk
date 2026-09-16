package config

import "os"

type Config struct {
	GRPCAddr       string
	DatabaseURL    string
	NATSURL        string
	TicketGRPCAddr string
}

func Load() Config {
	return Config{
		GRPCAddr:       getenv("ASSIGNMENT_GRPC_ADDR", ":50053"),
		DatabaseURL:    getenv("ASSIGNMENT_DATABASE_URL", "postgres://helpdesk:helpdesk@localhost:5432/assignment?sslmode=disable"),
		NATSURL:        getenv("NATS_URL", "nats://localhost:4222"),
		TicketGRPCAddr: getenv("TICKET_GRPC_ADDR", "localhost:50052"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
