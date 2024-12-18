package main

import (
	"context"
	"dobrify/botapp"
	"dobrify/dobry"
	"dobrify/internal/alog"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-telegram/bot"
)

func main() {
	var devMode bool
	if os.Getenv("DEV_MODE") == "1" {
		devMode = true
	}

	logger, close := alog.New("cron.log", devMode)
	defer close()

	botOpts := []bot.Option{
		bot.WithDebugHandler(func(format string, args ...any) {
			logger.Debug(fmt.Sprintf(format, args...))
		}),
		bot.WithErrorsHandler(func(err error) {
			logger.Error("bot error", alog.Error(err))
		}),
		bot.WithSkipGetMe(),
	}
	if devMode {
		botOpts = append(botOpts, bot.WithDebug())
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), botOpts...)
	if err != nil {
		logger.Error("failed to create bot", alog.Error(err))
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
			app.CheckPrizesAvailable(jobCtx, b, dobry.Glasses)
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
