package botapp

import (
	"context"
	"dobrify/crypter"
	"dobrify/dobry"
	"errors"
	"log/slog"
	"os"
)

var errUserNotFound = errors.New("user not found")

const filename = "app_state.bin"

type AppState struct {
	Pause       bool
	Users       map[string]*User
	NotifyUsers []string
}

type User struct {
	ChatID int64 `json:"cid"`
}

type App struct {
	secretKey     string
	adminUsername string
	dobryApp      *dobry.App
	state         *AppState
}

func NewApp(secretKey, adminUsername string) *App {
	var appState *AppState
	crypter.LoadFromFile(secretKey, filename, &appState)
	app := &App{
		secretKey:     secretKey,
		adminUsername: adminUsername,
		state:         appState,
	}
	if app.state == nil {
		app.state = app.initState()
	}
	app.fixupState()

	var username, password string
	if username = os.Getenv("DOBRY_USERNAME"); username == "" {
		slog.Error("DOBRY_USERNAME env variable must be provided")
		return app
	}
	if password = os.Getenv("DOBRY_PASSWORD"); password == "" {
		slog.Error("DOBRY_PASSWORD env variable must be provided")
		return app
	}

	app.dobryApp = dobry.NewApp(username, password, secretKey)
	return app
}

func (a *App) initState() *AppState {
	return &AppState{
		Pause:       false,
		Users:       make(map[string]*User),
		NotifyUsers: []string{a.adminUsername},
	}
}

func (a *App) fixupState() {
	if a.state.Users == nil {
		a.state.Users = make(map[string]*User)
	}
	if a.state.NotifyUsers == nil {
		a.state.NotifyUsers = []string{a.adminUsername}
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
	crypter.SaveToFile(a.secretKey, filename, a.state)
}
