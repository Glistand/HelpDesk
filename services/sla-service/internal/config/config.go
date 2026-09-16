package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCAddr         string
	DatabaseURL      string
	NATSURL          string
	RedisURL         string
	FirstResponse    time.Duration
	Resolve          time.Duration
	WarnRatio        float64
	PollInterval     time.Duration
}

func Load() Config {
	return Config{
		GRPCAddr:      getenv("SLA_GRPC_ADDR", ":50055"),
		DatabaseURL:   getenv("SLA_DATABASE_URL", "postgres://helpdesk:helpdesk@localhost:5432/sla?sslmode=disable"),
		NATSURL:       getenv("NATS_URL", "nats://localhost:4222"),
		RedisURL:      getenv("REDIS_URL", "redis://localhost:6379/0"),
		FirstResponse: durationEnv("SLA_FIRST_RESPONSE", 30*time.Minute),
		Resolve:       durationEnv("SLA_RESOLVE", 4*time.Hour),
		WarnRatio:     floatEnv("SLA_WARN_RATIO", 0.75),
		PollInterval:  durationEnv("SLA_POLL_INTERVAL", 250*time.Millisecond),
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
		return def
	}
	return d
}

func floatEnv(k string, def float64) float64 {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f <= 0 || f >= 1 {
		return def
	}
	return f
}
