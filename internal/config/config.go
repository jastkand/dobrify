package config

import (
	"errors"
	"fmt"
	"os"
)

var errMissingEnvVar = errors.New("missing env variable")

type Config struct {
	DevMode       bool
	BotToken      string
	DobryUsername string
	DobryPassword string
	AdminUsername string
	SecretKey     string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		DevMode:  os.Getenv("DEV_MODE") == "1",
		BotToken: os.Getenv("BOT_TOKEN"),
	}

	var err error
	if cfg.SecretKey, err = requireEnvVar("SECRET_KEY"); err != nil {
		return Config{}, err
	}
	if cfg.DobryUsername, err = requireEnvVar("DOBRY_USERNAME"); err != nil {
		return Config{}, err
	}
	if cfg.DobryPassword, err = requireEnvVar("DOBRY_PASSWORD"); err != nil {
		return Config{}, err
	}
	if cfg.AdminUsername, err = requireEnvVar("ADMIN_USERNAME"); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func requireEnvVar(name string) (string, error) {
	if v := os.Getenv(name); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("%w: %s", errMissingEnvVar, name)
}
