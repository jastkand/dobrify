package testutil

import (
	"testing"

	"github.com/go-telegram/bot"
)

func NewBotMock(t *testing.T) *bot.Bot {
	b, err := bot.New("token",
		bot.WithSkipGetMe(),
		bot.UseTestEnvironment(),
	)
	if err != nil {
		t.Fatalf("failed to create bot: %v", err)
	}
	return b
}
