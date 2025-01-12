package botapp

import (
	"context"
	"dobrify/dobry"
	"dobrify/internal/alog"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (a *App) getAvailablePrizes() (map[string]struct{}, error) {
	if a.dobryApp == nil {
		return nil, errDobryAppMissing
	}
	prizes, err := a.dobryApp.GetAvailablePrizes()
	if err != nil {
		return nil, err
	}
	prizesMap := make(map[string]struct{})
	for _, prize := range prizes {
		prizesMap[prize] = struct{}{}
	}
	return prizesMap, nil
}

func (a *App) CheckPrizesAvailable(ctx context.Context, b *bot.Bot) {
	if a.state.Pause {
		return
	}

	availablePrizes, err := a.getAvailablePrizes()
	if err != nil {
		slog.Error("failed to check prizes", alog.Error(err))
		return
	}
	if len(availablePrizes) == 0 {
		return
	}

	for _, user := range a.state.Users {
		if user.Pause {
			continue
		}
		shouldNotify := false
		text := "Доступны интересующие призы:\n"
		for _, prize := range dobry.Glasses { // TODO: put prizes user is interested in
			if _, exists := availablePrizes[prize]; exists {
				text += "\\+ " + dobry.PrizeName(prize) + "\n"
				shouldNotify = true
			}
		}
		if !shouldNotify {
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
