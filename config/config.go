package config

import "os"

type Config struct {
	BotToken          string
	DatabaseURL       string
	RenderExternalURL string
	Port              string
}

func Load() Config {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return Config{
		BotToken:          os.Getenv("BOT_TOKEN"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		RenderExternalURL: os.Getenv("RENDER_EXTERNAL_URL"),
		Port:              port,
	}
}