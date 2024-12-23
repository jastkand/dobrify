package botapp

import (
	"context"
	"dobrify/dobry"
	"dobrify/internal/alog"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) RegisterHandlers(ctx context.Context, b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, a.startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/check", bot.MatchTypeExact, a.checkHandler)
	// Admin commands
	b.RegisterHandler(bot.HandlerTypeMessageText, "/status", bot.MatchTypeExact, a.statusHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/sub", bot.MatchTypeExact, a.subscribeHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/sub ", bot.MatchTypePrefix, a.subscribeByUsernameHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/pause", bot.MatchTypeExact, a.pauseHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/resume", bot.MatchTypeExact, a.resumeHandler)

	b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Начать работу"},
			{Command: "check", Description: "Проверить доступные призы"},
		},
	})

	if adminUser, exists := a.state.Users[a.cfg.AdminUsername]; exists {
		b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
			Commands: []models.BotCommand{
				{Command: "start", Description: "Начать работу"},
				{Command: "check", Description: "Проверить доступные призы"},
				{Command: "status", Description: "Показать статус"},
				{Command: "sub", Description: "Подписать пользователя"},
				{Command: "pause", Description: "Приостановить работу"},
				{Command: "resume", Description: "Возобновить работу"},
			},
			Scope: &models.BotCommandScopeChat{ChatID: adminUser.ChatID},
		})
	}
}

func (a *App) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.Username == a.cfg.AdminUsername && a.subUserState == subUserStateUsername {
		if ok := a.handleUserSubscribe(ctx, b, update, update.Message.Text); ok {
			a.subUserState = subUserStateNone
		}
		return
	}
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

func (a *App) statusHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("status", "username", update.Message.From.Username)
	if update.Message.From.Username != a.cfg.AdminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	a.subUserState = subUserStateNone
	var text string
	if a.state.Pause {
		text = "Отдыхаю."
	} else {
		text = "Работаю."
	}
	if len(a.state.NotifyUsers) > 0 {
		usersNames := make([]string, 0, len(a.state.NotifyUsers))
		for _, uname := range a.state.NotifyUsers {
			usersNames = append(usersNames, "@"+uname)
		}
		text += "\n\nРазрешено оповещать: " + strings.Join(usersNames, ", ")
	} else {
		text += "\n\nНикому ничего не скажу."
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}

func (a *App) subscribeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("subscribe", "username", update.Message.From.Username)
	if update.Message.From.Username != a.cfg.AdminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	a.subUserState = subUserStateUsername
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Введи юзернейм пользователя для подписки",
	})
}

func (a *App) subscribeByUsernameHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("subscribe by username", "username", update.Message.From.Username, "text", update.Message.Text)
	if update.Message.From.Username != a.cfg.AdminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	a.subUserState = subUserStateNone
	text := update.Message.Text
	if len(text) <= 5 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "А кого подписывать? Юзернейм надо.",
		})
		return
	}
	a.handleUserSubscribe(ctx, b, update, text[5:])
}

func (a *App) handleUserSubscribe(ctx context.Context, b *bot.Bot, update *models.Update, username string) bool {
	slog.Debug("handle user subscribe", "username", update.Message.From.Username, "sub", username)
	_, err := a.subscribeUser(ctx, username)
	if err != nil {
		slog.Error("failed to subscribe user", alog.Error(err), "username", username)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пользаватель @" + username + " не зарегистрирован.",
		})
		return false
	}
	slog.Debug("subscribed user", "username", username)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Пользаватель @" + username + " подписан.",
	})
	return true
}

func (a *App) pauseHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("pause", "username", update.Message.From.Username)
	if update.Message.From.Username != a.cfg.AdminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	a.subUserState = subUserStateNone
	a.pause(ctx)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Отдыхаю.",
	})
}

func (a *App) resumeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Debug("resume", "username", update.Message.From.Username)
	if update.Message.From.Username != a.cfg.AdminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	a.subUserState = subUserStateNone
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
	prizes, err := hasWantedPrizes(a, dobry.AllPrizes)
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
		for _, prize := range prizes {
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

func hasWantedPrizes(a *App, wanted []string) ([]string, error) {
	if a.dobryApp == nil {
		return nil, errDobryAppMissing
	}
	prizes, err := a.dobryApp.HasWantedPrizes(wanted)
	if err != nil {
		return nil, err
	}
	return prizes, nil
}
