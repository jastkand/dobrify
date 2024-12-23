package botapp

import "strings"

func normalizeUsername(username string) string {
	return strings.Trim(strings.ToLower(username), "@")
}
