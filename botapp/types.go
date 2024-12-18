package botapp

import "errors"

var (
	errUserNotFound    = errors.New("user not found")
	errDobryAppMissing = errors.New("dobry app missing")
	errAppPaused       = errors.New("app is paused")
)
