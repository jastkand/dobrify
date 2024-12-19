package main

import (
	"dobrify/dobry"
	"dobrify/internal/alog"
	"dobrify/internal/config"
)

func main() {
	logger, close := alog.New("check.log", config.IsDevStage())
	defer close()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", alog.Error(err))
		return
	}

	app := dobry.NewApp(cfg)
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
