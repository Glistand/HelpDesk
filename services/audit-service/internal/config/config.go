package config

import "os"

type Config struct {
	GRPCAddr    string
	DatabaseURL string
	NATSURL     string
}

func Load() Config {
	return Config{
		GRPCAddr:    getenv("AUDIT_GRPC_ADDR", ":50054"),
		DatabaseURL: getenv("AUDIT_DATABASE_URL", "postgres://helpdesk:helpdesk@localhost:5432/audit?sslmode=disable"),
		NATSURL:     getenv("NATS_URL", "nats://localhost:4222"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
