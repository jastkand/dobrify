package botapp

import (
	"context"
	"dobrify/dobry"
	"errors"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/hello", bot.MatchTypeExact, a.helloHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/status", bot.MatchTypeExact, a.statusHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/sub ", bot.MatchTypePrefix, a.subscribeHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/pause", bot.MatchTypeExact, a.pauseHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/resume", bot.MatchTypeExact, a.resumeHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/check", bot.MatchTypeExact, a.checkHandler)
}

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "You're probably looking for something else 🤔",
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

func (a *App) helloHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	a.addUser(ctx, update.Message.From.Username, update.Message.Chat.ID)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Привет, *" + bot.EscapeMarkdown(update.Message.From.FirstName) + "*",
		ParseMode: models.ParseModeMarkdown,
	})
}

func (a *App) statusHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.Username != a.adminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
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
		text += "\n\nБуду оповещать: " + strings.Join(usersNames, ", ")
	} else {
		text += "\n\nНикому ничего не скажу."
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}

func (a *App) subscribeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.Username != a.adminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	text := update.Message.Text
	if len(text) <= 5 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "А кого подписывать? Юзернейм надо.",
		})
		return
	}
	username := text[5:]
	subbed, err := a.subscribeUser(ctx, username)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "User @" + username + " is not registered.",
		})
		return
	}
	if subbed {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "User @" + username + " subscribed.",
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "User @" + username + " is already subscribed.",
		})
	}
}

func (a *App) pauseHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.Username != a.adminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	a.pause(ctx)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Отдыхаю.",
	})
}

func (a *App) resumeHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.Username != a.adminUsername {
		a.adminGuard(ctx, b, update)
		return
	}
	if a.state.Pause == false {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Я и так работаю.",
		})
		return
	}
	a.resume(ctx)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Готов к роботе.",
	})
}

func (a *App) checkHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	canCheck := false
	for _, uname := range a.state.NotifyUsers {
		if uname == update.Message.From.Username {
			canCheck = true
			break
		}
	}
	if !canCheck {
		a.adminGuard(ctx, b, update)
		return
	}
	prizes, err := hasWantedPrizes(a, dobry.AllPrizes)
	if err != nil {
		message := "Ошибка при запросе призов."
		if errors.Is(err, errAppPaused) {
			message = "Отдыхаю. Не могу проверить призы."
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   message,
		})
		return
	}
	if len(prizes) > 0 {
		text := "Доступны интересующие призы:\n"
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
			Text:   "Интересующих призов нет в наличии.",
		})
	}
}

func hasWantedPrizes(a *App, wanted []string) ([]string, error) {
	if a.state.Pause {
		return nil, errAppPaused
	}
	if a.dobryApp == nil {
		return nil, errDobryAppMissing
	}
	prizes, err := a.dobryApp.HasWantedPrizes(wanted)
	if err != nil {
		slog.Error("failed to check for wanted prizes", "error", err.Error())
		return nil, err
	}
	return prizes, nil
}
