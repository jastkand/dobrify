package botify

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

func Run(cfg config.Config) {
	logger := alog.New(cfg.IsDev())

	app := botapp.NewApp(cfg)

	botOpts := []bot.Option{
		bot.WithDebugHandler(func(format string, args ...any) {
			logger.Debug(fmt.Sprintf(format, args...))
		}),
		bot.WithErrorsHandler(func(err error) {
			logger.Error("bot error", alog.Error(err))
		}),
		bot.WithDefaultHandler(app.DefaultHandler),
	}
	if cfg.IsDev() {
		botOpts = append(botOpts, bot.WithDebug())
	}

	b, err := bot.New(cfg.BotToken, botOpts...)
	if err != nil {
		logger.Error("failed to create bot", alog.Error(err))
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app.RegisterHandlers(ctx, b)
	b.Start(ctx)
}
