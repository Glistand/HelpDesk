package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCAddr       string
	DatabaseURL    string
	NATSURL        string
	OutboxInterval time.Duration
}

func Load() Config {
	return Config{
		GRPCAddr:       getenv("TICKET_GRPC_ADDR", ":50052"),
		DatabaseURL:    getenv("TICKET_DATABASE_URL", "postgres://helpdesk:helpdesk@localhost:5432/ticket?sslmode=disable"),
		NATSURL:        getenv("NATS_URL", "nats://localhost:4222"),
		OutboxInterval: durationEnv("OUTBOX_POLL_INTERVAL", 500*time.Millisecond),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func durationEnv(k string, def time.Duration) time.Duration {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		if ms, err2 := strconv.Atoi(v); err2 == nil {
			return time.Duration(ms) * time.Millisecond
		}
		return def
	}
	return d
}
