package main

import (
	"context"
	"dobrify/botapp"
	"dobrify/internal/alog"
	"dobrify/internal/config"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

func main() {
	devMode := os.Getenv("DEV_MODE") == "1"

	logger, close := alog.New("bot.log", devMode)
	defer close()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", alog.Error(err))
		return
	}

	botOpts := []bot.Option{
		bot.WithDebugHandler(func(format string, args ...any) {
			logger.Debug(fmt.Sprintf(format, args...))
		}),
		bot.WithErrorsHandler(func(err error) {
			logger.Error("bot error", alog.Error(err))
		}),
		bot.WithDefaultHandler(botapp.DefaultHandler),
	}
	if devMode {
		botOpts = append(botOpts, bot.WithDebug())
	}

	b, err := bot.New(cfg.BotToken, botOpts...)
	if err != nil {
		logger.Error("failed to create bot", alog.Error(err))
		return
	}

	app := botapp.NewApp(cfg)
	app.RegisterHandlers(b)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	b.Start(ctx)
}
