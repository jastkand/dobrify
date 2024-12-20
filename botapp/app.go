package botapp

import (
	"context"
	"dobrify/dobry"
	"dobrify/internal/alog"
	"dobrify/internal/config"
	"dobrify/storage"
	"log/slog"
	"time"
)

const storeFilename = "app_state.json"

type AppState struct {
	Pause       bool
	Users       map[string]*User
	NotifyUsers []string
	UpdatedAt   time.Time
}

type User struct {
	ChatID int64 `json:"cid"`
}

type App struct {
	cfg      config.Config
	store    storage.Storage
	dobryApp *dobry.App
	state    *AppState
}

func NewApp(cfg config.Config) *App {
	store := storage.NewPlainStore()

	var appState AppState
	if err := store.LoadFromFile(storeFilename, &appState); err != nil {
		slog.Error("failed to load json state", alog.Error(err), "filename", storeFilename)
	}

	app := &App{
		cfg:   cfg,
		store: store,
		state: &appState,
	}
	if app.state == nil {
		app.state = app.initState()
	}
	app.fixupState()

	app.dobryApp = dobry.NewApp(cfg)
	return app
}

func (a *App) initState() *AppState {
	return &AppState{
		Pause:       false,
		Users:       make(map[string]*User),
		NotifyUsers: []string{a.cfg.AdminUsername},
	}
}

func (a *App) fixupState() {
	if a.state.Users == nil {
		a.state.Users = make(map[string]*User)
	}
	if a.state.NotifyUsers == nil {
		a.state.NotifyUsers = []string{a.cfg.AdminUsername}
	}
	notifyUsers := make([]string, 0, len(a.state.NotifyUsers))
	for _, uname := range a.state.NotifyUsers {
		if a.state.Users[uname] != nil {
			notifyUsers = append(notifyUsers, uname)
		}
	}
	a.state.NotifyUsers = notifyUsers
}

func (a *App) pause(ctx context.Context) {
	slog.Info("pausing")
	a.state.Pause = true
	go a.saveState(ctx)
}

func (a *App) resume(ctx context.Context) {
	slog.Info("resuming")
	a.state.Pause = false
	go a.saveState(ctx)
}

func (a *App) addUser(ctx context.Context, username string, chatID int64) bool {
	if a.state.Users == nil {
		a.state.Users = make(map[string]*User)
	}
	if _, ok := a.state.Users[username]; ok {
		return false
	}
	slog.Info("adding user", "username", username, "chatID", chatID)
	a.state.Users[username] = &User{ChatID: chatID}
	go a.saveState(ctx)
	return true
}

func (a *App) subscribeUser(ctx context.Context, username string) (bool, error) {
	for _, uname := range a.state.NotifyUsers {
		if a.state.Users[uname] == nil {
			slog.Warn("user not found", "username", uname)
			return false, errUserNotFound
		}
		if uname == username {
			return false, nil
		}
	}
	slog.Info("subscribing user", "username", username)
	a.state.NotifyUsers = append(a.state.NotifyUsers, username)
	go a.saveState(ctx)
	return true, nil
}

func (a *App) saveState(ctx context.Context) {
	if ctx.Err() != nil {
		slog.Debug("context is done, not saving state")
		return
	}
	a.state.UpdatedAt = time.Now()
	if err := a.store.SaveToFile(storeFilename, a.state); err != nil {
		slog.Error("failed to save json state", alog.Error(err), "filename", storeFilename)
	}
}
