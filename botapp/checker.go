package botapp

import (
	"context"
	"dobrify/dobry"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) CheckPrizesAvailable(ctx context.Context, b *bot.Bot, wanted []string) {
	prizes, err := hasWantedPrizes(a, wanted)
	if err != nil || len(prizes) == 0 {
		return
	}

	text := "Доступны интересующие призы:\n"
	for _, prize := range prizes {
		text += "\\- " + dobry.PrizeName(prize) + "\n"
	}
	for _, username := range a.state.NotifyUsers {
		if user, ok := a.state.Users[username]; ok {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    user.ChatID,
				Text:      text,
				ParseMode: models.ParseModeMarkdown,
			})
			if err != nil {
				slog.Error("failed to send message", "error", err.Error())
			}
		}
	}
}
