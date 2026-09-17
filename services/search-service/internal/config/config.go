package config

import "os"

type Config struct {
	GRPCAddr       string
	NATSURL        string
	TicketGRPCAddr string
	MeiliHost      string
	MeiliAPIKey    string
	MeiliIndex     string
}

func Load() Config {
	return Config{
		GRPCAddr:       getenv("SEARCH_GRPC_ADDR", ":50057"),
		NATSURL:        getenv("NATS_URL", "nats://localhost:4222"),
		TicketGRPCAddr: getenv("TICKET_GRPC_ADDR", "localhost:50052"),
		MeiliHost:      getenv("MEILI_HOST", "http://localhost:7700"),
		MeiliAPIKey:    getenv("MEILI_MASTER_KEY", "helpdesk-dev-key"),
		MeiliIndex:     getenv("MEILI_INDEX", "tickets"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
