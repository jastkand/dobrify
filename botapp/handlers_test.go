package botapp

import (
	"context"
	"dobrify/storage"
	"dobrify/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

func newApp(state *AppState) *App {
	app := &App{
		store: storage.NewInMemoryStore(nil),
		state: state,
	}
	if app.state == nil {
		app.state = app.initState()
	}
	return app
}

func TestApp_pauseHandler(t *testing.T) {
	t.Parallel()

	t.Run("pauses the user notifications", func(t *testing.T) {
		t.Parallel()
		app := newApp(&AppState{
			Users: map[string]*User{
				"test": {Pause: false},
			},
		})
		app.pauseHandler(context.Background(), testutil.NewBotMock(t), &models.Update{
			Message: &models.Message{
				From: &models.User{
					Username: "test",
				},
			},
		})
		if !app.state.Users["test"].Pause {
			t.Error("expected user to be paused")
		}
	})
}

func TestApp_resumeHandler(t *testing.T) {
	t.Parallel()

	t.Run("resumes the user notifications", func(t *testing.T) {
		t.Parallel()
		app := newApp(&AppState{
			Users: map[string]*User{
				"test": {Pause: true},
			},
		})
		app.resumeHandler(context.Background(), testutil.NewBotMock(t), &models.Update{
			Message: &models.Message{
				From: &models.User{
					Username: "test",
				},
			},
		})
		if app.state.Users["test"].Pause {
			t.Error("expected user to be resumed")
		}
	})
}
