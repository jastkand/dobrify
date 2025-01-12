package botapp

import (
	"context"
	"dobrify/dobry"
	"dobrify/internal/alog"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) RegisterHandlers(ctx context.Context, b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, a.startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/check", bot.MatchTypeExact, a.checkHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/pause", bot.MatchTypeExact, a.pauseHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/resume", bot.MatchTypeExact, a.resumeHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stop", bot.MatchTypeExact, a.stopHandler)
	// Admin commands
	b.RegisterHandler(bot.HandlerTypeMessageText, "/pause_all", bot.MatchTypeExact, a.pauseAllHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/resume_all", bot.MatchTypeExact, a.resumeAllHandler)

	b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Начать работу"},
			{Command: "check", Description: "Проверить доступные призы"},
			{Command: "pause", Description: "Приостановить отправку уведомлений"},
			{Command: "resume", Description: "Возобновить отправку уведомлений"},
			{Command: "stop", Description: "Закончить работу"},
		},
	})

	if adminUser, exists := a.state.Users[a.cfg.AdminUsername]; exists {
		b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
			Commands: []models.BotCommand{
				{Command: "check", Description: "Проверить доступные призы"},
				{Command: "pause", Description: "Приостановить отправку уведомлений"},
				{Command: "resume", Description: "Возобновить отправку уведомлений"},
				{Command: "pause_all", Description: "Приостановить работу"},
				{Command: "resume_all", Description: "Возобновить работу"},
			},
			Scope: &models.BotCommandScopeChat{ChatID: adminUser.ChatID},
		})
	}
}

func (a *App) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ничего не понятно 🤔",
		})
	}
}

func (a *App) adminGuard(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info("admin guard", "username", update.Message.From.Username)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   `Сюда нельзя ¯\_(ツ)_/¯`,
	})
}

func (a *App) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("start", "username", update.Message.From.Username)
	a.addUser(ctx, update.Message.From.Username, update.Message.Chat.ID)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Привет, *" + bot.EscapeMarkdown(update.Message.From.FirstName) + "*",
		ParseMode: models.ParseModeMarkdown,
	})
}

func (a *App) stopHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("stop", "username", update.Message.From.Username)
	a.removeUser(ctx, update.Message.From.Username)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Пока, *" + bot.EscapeMarkdown(update.Message.From.FirstName) + "*",
		ParseMode: models.ParseModeMarkdown,
	})
}

func (a *App) pauseHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("pause", "username", update.Message.From.Username)
	a.pauseUser(ctx, update.Message.From.Username)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Отдыхаю.",
	})
}

func (a *App) resumeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("resume", "username", update.Message.From.Username)
	a.resumeUser(ctx, update.Message.From.Username)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Готов к роботе.",
	})
}

func (a *App) pauseAllHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("resumeAll", "username", update.Message.From.Username)
	if update.Message.From.Username != a.cfg.AdminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	a.pause(ctx)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Отдыхаю.",
	})
}

func (a *App) resumeAllHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("resumeAll", "username", update.Message.From.Username)
	if update.Message.From.Username != a.cfg.AdminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	if a.state.Pause != false {
		a.resume(ctx)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Готов к роботе.",
	})
}

func (a *App) checkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("check", "username", update.Message.From.Username)
	prizes, err := a.getAvailablePrizes()
	if err != nil {
		slog.Error("failed to check prizes", alog.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при запросе призов.",
		})
		return
	}
	if len(prizes) > 0 {
		text := "Список доступных призов:\n"
		for prize := range prizes {
			text += "- " + dobry.PrizeName(prize) + "\n"
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text,
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Нет доступных призов.",
		})
	}
}
