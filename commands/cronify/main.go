package cronify

import (
	"context"
	"dobrify/botapp"
	"dobrify/internal/alog"
	"dobrify/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-telegram/bot"
)

func Run(cfg config.Config) {
	logger := alog.New(config.IsDevStage())

	botOpts := []bot.Option{
		bot.WithDebugHandler(func(format string, args ...any) {
			logger.Debug(fmt.Sprintf(format, args...))
		}),
		bot.WithErrorsHandler(func(err error) {
			logger.Error("bot error", alog.Error(err))
		}),
		bot.WithSkipGetMe(),
	}
	if cfg.IsDev() {
		botOpts = append(botOpts, bot.WithDebug())
	}

	b, err := bot.New(cfg.BotToken, botOpts...)
	if err != nil {
		logger.Error("failed to create bot", alog.Error(err))
		return
	}

	jobCtx, jobCancel := context.WithCancel(context.Background())
	defer jobCancel()

	slog.Debug("starting scheduler")
	s, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("failed to create scheduler", alog.Error(err))
		return
	}
	_, err = s.NewJob(
		gocron.CronJob("*/10 * * * *", false),
		gocron.NewTask(func() {
			app := botapp.NewApp(cfg)
			app.CheckPrizesAvailable(jobCtx, b)
		}),
	)
	if err != nil {
		slog.Error("failed to create job", alog.Error(err))
		return
	}

	s.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	done := make(chan struct{})

	go func() {
		<-quit
		slog.Debug("shutdown signal received")

		jobCancel()
		if err = s.Shutdown(); err != nil {
			slog.Error("failed to shutdown scheduler", alog.Error(err))
		}

		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCancel()

		slog.Debug("waiting for shutdown to complete")
		<-timeoutCtx.Done()

		done <- struct{}{}
	}()

	<-done
}
