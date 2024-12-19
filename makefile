
bot.dev: .env
	go run cmd/bot/main.go

check.dev: .env
	go run cmd/check/main.go

cron.dev: .env
	go run cmd/cron/main.go
