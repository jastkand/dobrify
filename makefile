ifeq (runner,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

runner:
	@go run cmd/runner/main.go $(RUN_ARGS)
	
bot.dev:
	@make runner bot

cron.dev:
	@make runner cron

build:
	@echo "Building..."
	@STAGE=prod go build -o dist/dobrify cmd/runner/main.go

build.linux:
	@echo "Building for linux..."
	STAGE=prod GOOS=linux GOARCH=amd64 go build -o dist/dobrify-linux cmd/runner/main.go
	@echo "Build complete!"

include .env.deploy
deploy: .env.deploy build.linux
	@echo "Deploying..."	
	@scp -rp ./dist/dobrify-linux $(DEPLOY_USER)@$(DEPLOY_HOST):$(DEPLOY_PATH)
	@ssh $(DEPLOY_USER)@$(DEPLOY_HOST) "sh $(DEPLOY_PATH)/after_deploy.sh"
	@echo "Deploy complete!"