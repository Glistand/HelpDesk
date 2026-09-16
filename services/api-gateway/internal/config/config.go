package config

import "os"

type Config struct {
	HTTPAddr       string
	AuthGRPCAddr   string
	TicketGRPCAddr string
}

func Load() Config {
	return Config{
		HTTPAddr:       getenv("GATEWAY_HTTP_ADDR", ":8080"),
		AuthGRPCAddr:   getenv("AUTH_GRPC_ADDR", "localhost:50051"),
		TicketGRPCAddr: getenv("TICKET_GRPC_ADDR", "localhost:50052"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
