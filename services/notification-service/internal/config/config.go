package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL   string
	NATSURL       string
	FailFirstN    int // for retry demos: fail first N attempts per process
}

func Load() Config {
	return Config{
		DatabaseURL: getenv("NOTIFICATION_DATABASE_URL", "postgres://helpdesk:helpdesk@localhost:5432/notification?sslmode=disable"),
		NATSURL:     getenv("NATS_URL", "nats://localhost:4222"),
		FailFirstN:  intEnv("NOTIFY_FAIL_FIRST_N", 0),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func intEnv(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
