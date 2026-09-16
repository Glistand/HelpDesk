package config

import (
	"os"
	"time"
)

type Config struct {
	GRPCAddr  string
	JWTSecret string
	TokenTTL  time.Duration
}

func Load() Config {
	return Config{
		GRPCAddr:  getenv("AUTH_GRPC_ADDR", ":50051"),
		JWTSecret: getenv("JWT_SECRET", "helpdesk-dev-secret-change-me"),
		TokenTTL:  24 * time.Hour,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
