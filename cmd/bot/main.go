package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
		bot.WithDefaultHandler(handler),
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

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Printf("%+v\n", ctx)
	fmt.Printf("%+v\n", b)
	fmt.Printf("%+v\n", update)
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   update.Message.Text,
		})
	}
}
