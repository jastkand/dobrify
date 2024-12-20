package main

import (
	"dobrify/commands/botify"
	"dobrify/commands/cronify"
	"dobrify/commands/manual"
	"dobrify/internal/alog"
	"dobrify/internal/config"
	"fmt"
	"os"
)

type Runnable func(config.Config)

type Command struct {
	Name        string
	Description string
	Run         Runnable
}

var commands = map[string]*Command{
	"cron": {
		Name:        "cron",
		Description: "Run cron job",
		Run:         cronify.Run,
	},
	"bot": {
		Name:        "bot",
		Description: "Run bot",
		Run:         botify.Run,
	},
	"manual": {
		Name:        "manual",
		Description: "Manually check for prizes availability",
		Run:         manual.Run,
	},
	"help": {
		Name:        "help",
		Description: "Show available commands",
	},
}

func main() {
	logger := alog.New(config.IsDevStage())

	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("Usage: dobrify <command>")
		fmt.Println()
		listAvailableCommands()
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", alog.Error(err))
		return
	}

	commandName := os.Args[1]

	if commandName == "help" {
		listAvailableCommands()
		return
	}

	command, ok := commands[commandName]
	if !ok {
		fmt.Println("Unknown command:", commandName)
		fmt.Println()
		listAvailableCommands()
		return
	}

	command.Run(cfg)
}

func listAvailableCommands() {
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("  %s: %s\n", cmd.Name, cmd.Description)
	}
}
