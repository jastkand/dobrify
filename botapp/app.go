package botapp

import (
	"context"
	"dobrify/dobry"
	"dobrify/internal/alog"
	"dobrify/internal/config"
	"dobrify/storage"
	"log/slog"
	"slices"
	"time"
)

type subUserState int

const (
	subUserStateNone subUserState = iota
	subUserStateUsername
)

type AppState struct {
	Pause       bool             `json:"pause"`
	Users       map[string]*User `json:"users"`
	NotifyUsers []string         `json:"notify_users"`
	Version     int64            `json:"v"`
}

type User struct {
	ChatID int64 `json:"cid"`
	Pause  bool  `json:"pause"`
}

type App struct {
	cfg          config.Config
	store        storage.Storage
	dobryApp     *dobry.App
	state        *AppState
	subUserState subUserState
}

func NewApp(cfg config.Config) *App {
	store := storage.NewJSONStore()

	var appState AppState
	if err := store.LoadFromFile(cfg.StorePath, &appState); err != nil {
		slog.Error("failed to load json state", alog.Error(err), "filename", cfg.StorePath)
	}

	app := &App{
		cfg:          cfg,
		store:        store,
		state:        &appState,
		subUserState: subUserStateNone,
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
	slog.Debug("pause")
	a.state.Pause = true
	go a.saveState(ctx)
}

func (a *App) resume(ctx context.Context) {
	slog.Debug("resume")
	a.state.Pause = false
	go a.saveState(ctx)
}

func (a *App) pauseUser(ctx context.Context, username string) {
	slog.Debug("pause user", "username", username)
	normalized := normalizeUsername(username)
	if user, exists := a.state.Users[normalized]; exists {
		user.Pause = true
		go a.saveState(ctx)
	}
}

func (a *App) resumeUser(ctx context.Context, username string) {
	slog.Debug("resume user", "username", username)
	normalized := normalizeUsername(username)
	if user, exists := a.state.Users[normalized]; exists {
		user.Pause = false
		go a.saveState(ctx)
	}
}

func (a *App) addUser(ctx context.Context, username string, chatID int64) bool {
	slog.Debug("add user", "username", username, "chatID", chatID)
	if a.state.Users == nil {
		a.state.Users = make(map[string]*User)
	}
	normalized := normalizeUsername(username)
	if _, exists := a.state.Users[normalized]; exists {
		return false
	}
	slog.Info("adding user", "username", username, "chatID", chatID)
	a.state.Users[normalized] = &User{ChatID: chatID}
	go a.saveState(ctx)
	return true
}

func (a *App) removeUser(ctx context.Context, username string) bool {
	slog.Debug("remove user", "username", username)
	normalized := normalizeUsername(username)
	a.state.NotifyUsers = slices.DeleteFunc(a.state.NotifyUsers, func(el string) bool {
		return el == normalized
	})
	if _, exists := a.state.Users[normalized]; exists {
		delete(a.state.Users, normalized)
	}
	go a.saveState(ctx)
	return true
}

func (a *App) subscribeUser(ctx context.Context, username string) (bool, error) {
	slog.Debug("subscribe user", "username", username)
	normalized := normalizeUsername(username)
	if _, exists := a.state.Users[normalized]; !exists {
		slog.Warn("user not found", "username", username)
		return false, errUserNotFound
	}
	for _, uname := range a.state.NotifyUsers {
		if uname == normalized {
			slog.Warn("user already subscribed", "username", username)
			return false, nil
		}
	}
	slog.Info("subscribing user", "username", username)
	a.state.NotifyUsers = append(a.state.NotifyUsers, normalized)
	go a.saveState(ctx)
	return true, nil
}

func (a *App) saveState(ctx context.Context) {
	slog.Debug("save state")
	if ctx.Err() != nil {
		slog.Debug("context is done, not saving state")
		return
	}
	a.state.Version = time.Now().UnixMilli()
	if err := a.store.SaveToFile(a.cfg.StorePath, a.state); err != nil {
		slog.Error("failed to save json state", alog.Error(err), "filename", a.cfg.StorePath)
	}
}
