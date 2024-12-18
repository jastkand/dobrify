package main

import (
	"context"
	"dobrify/botapp"
	"dobrify/dobry"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

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

	botOpts := []bot.Option{
		bot.WithSkipGetMe(),
	}
	if devMode {
		botOpts = append(botOpts, bot.WithDebug())
		botOpts = append(botOpts, bot.WithDebugHandler(func(format string, args ...any) {
			logger.Info(fmt.Sprintf(format, args...))
		}))
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), botOpts...)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create bot: %v", err.Error()))
		return
	}

	var secretKey, adminUsername string
	if secretKey = os.Getenv("SECRET_KEY"); secretKey == "" {
		logger.Error("SECRET_KEY env variable must be provided")
		return
	}
	if adminUsername = os.Getenv("ADMIN_USERNAME"); adminUsername == "" {
		logger.Error("ADMIN_USERNAME env variable must be provided")
		return
	}
	app := botapp.NewApp(secretKey, adminUsername)
	app.CheckPrizesAvailable(ctx, b, dobry.Glasses)
}
