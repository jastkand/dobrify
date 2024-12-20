package config

import (
	"dobrify/internal/alog"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Stage         string `env:"STAGE"`
	BotToken      string `env:"BOT_TOKEN"`
	DobryUsername string `env:"DOBRY_USERNAME,required"`
	DobryPassword string `env:"DOBRY_PASSWORD,required"`
	AdminUsername string `env:"ADMIN_USERNAME,required"`
	SecretKey     string `env:"SECRET_KEY,required"`
}

func (c Config) IsDev() bool {
	return c.Stage == "dev"
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env file", alog.Error(err))
		slog.Debug(".env file is missing, using environment variables")
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Stage == "" {
		cfg.Stage = "dev"
	}

	return cfg, nil
}

func IsDevStage() bool {
	envStage := os.Getenv("STAGE")
	return envStage == "" || envStage == "dev"
}
