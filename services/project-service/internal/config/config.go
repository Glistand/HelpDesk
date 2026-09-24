package config

import "os"

type BootstrapProfile struct {
	Name             string
	Description      string
	Domain           string
	WebsiteURL       string
	SupportEmail     string
	SupportPhone     string
	Timezone         string
	SiteKey          string
	BotInstructions  string
	AutoCreateTicket bool
}

type Config struct {
	GRPCAddr    string
	DatabaseURL string
	Bootstrap   BootstrapProfile
}

func Load() Config {
	return Config{
		GRPCAddr:    getenv("PROJECT_GRPC_ADDR", ":50059"),
		DatabaseURL: getenv("PROJECT_DATABASE_URL", "postgres://helpdesk:helpdesk@localhost:5432/project?sslmode=disable"),
		Bootstrap: BootstrapProfile{
			Name:             getenv("PROJECT_NAME", "SalonPro"),
			Description:      getenv("PROJECT_DESCRIPTION", "Платформа для управления салонами красоты: запись, оплата и рабочие места администраторов и мастеров."),
			Domain:           getenv("PROJECT_DOMAIN", "beauty-salon-management"),
			WebsiteURL:       os.Getenv("PROJECT_WEBSITE_URL"),
			SupportEmail:     os.Getenv("PROJECT_SUPPORT_EMAIL"),
			SupportPhone:     os.Getenv("PROJECT_SUPPORT_PHONE"),
			Timezone:         getenv("PROJECT_TIMEZONE", "Europe/Moscow"),
			SiteKey:          getenv("WIDGET_SITE_KEY", "demo-site"),
			BotInstructions:  getenv("PROJECT_BOT_INSTRUCTIONS", "Помогай клиентам записаться на услуги салона, узнать стоимость, изменить или отменить запись. Если уверенного ответа нет, передай обращение менеджеру."),
			AutoCreateTicket: getenv("PROJECT_AUTO_CREATE_TICKET", "true") != "false",
		},
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
