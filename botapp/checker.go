package botapp

import (
	"context"
	"dobrify/dobry"
	"dobrify/internal/alog"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) CheckPrizesAvailable(ctx context.Context, b *bot.Bot, wanted []string) {
	if a.state.Pause {
		return
	}

	prizes, err := hasWantedPrizes(a, wanted)
	if err != nil {
		slog.Error("failed to check prizes", alog.Error(err))
		return
	}
	if len(prizes) == 0 {
		return
	}

	text := "Доступны интересующие призы:\n"
	for _, prize := range prizes {
		text += "\\- " + dobry.PrizeName(prize) + "\n"
	}
	for _, username := range a.state.NotifyUsers {
		if user, ok := a.state.Users[username]; ok {
			if user.Pause {
				continue
			}
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    user.ChatID,
				Text:      text,
				ParseMode: models.ParseModeMarkdown,
			})
			if err != nil {
				slog.Error("failed to send message", alog.Error(err))
			}
		}
	}
}
