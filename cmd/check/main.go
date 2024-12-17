package main

import (
	"dobrify/dobry"
	"log/slog"
	"os"
)

func main() {
	var devMode bool
	if os.Getenv("DEV_MODE") == "1" {
		devMode = true
	}
	loggerOpts := &slog.HandlerOptions{}
	if devMode {
		loggerOpts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, loggerOpts))
	slog.SetDefault(logger)

	var username, password, secretKey string
	if username = os.Getenv("DOBRY_USERNAME"); username == "" {
		logger.Error("DOBRY_USERNAME env variable must be provided")
		return
	}
	if password = os.Getenv("DOBRY_PASSWORD"); password == "" {
		logger.Error("DOBRY_PASSWORD env variable must be provided")
		return
	}
	if secretKey = os.Getenv("SECRET_KEY"); secretKey == "" {
		logger.Error("SECRET_KEY env variable must be provided")
		return
	}

	app := dobry.NewApp(username, password, secretKey)
	prizes, err := app.HasWantedPrizes(dobry.Glasses)
	if err != nil {
		logger.Error("failed to check for wanted prizes", "error", err.Error())
		return
	}

	if len(prizes) > 0 {
		logger.Info("you have a wanted prizes", "prizes", prizes)
	} else {
		logger.Info("you don't have a wanted prize")
	}
}
