package config

import "os"

type Config struct {
	GRPCAddr          string
	DatabaseURL       string
	NATSURL           string
	OpenRouterAPIKey  string
	OpenRouterModel   string
	OpenRouterBaseURL string
	SiteKey           string
}

func Load() Config {
	return Config{
		GRPCAddr:          getenv("CONVERSATION_GRPC_ADDR", ":50058"),
		DatabaseURL:       getenv("CONVERSATION_DATABASE_URL", "postgres://helpdesk:helpdesk@localhost:5432/conversation?sslmode=disable"),
		NATSURL:           getenv("NATS_URL", "nats://localhost:4222"),
		OpenRouterAPIKey:  os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel:   getenv("OPENROUTER_MODEL", "inclusionai/ling-3.0-flash-vl:free"),
		OpenRouterBaseURL: getenv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
		SiteKey:           getenv("WIDGET_SITE_KEY", "demo-site"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
