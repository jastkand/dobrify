package manual

import (
	"dobrify/dobry"
	"dobrify/internal/config"
	"fmt"
	"strings"
)

func Run(cfg config.Config) {
	app := dobry.NewApp(cfg)
	prizes, err := app.HasWantedPrizes(dobry.Glasses)
	if err != nil {
		fmt.Printf("failed to check for wanted prizes: %v\n", err)
		return
	}

	if len(prizes) > 0 {
		fmt.Printf("Prizes are available: %s\n", strings.Join(prizes, ", "))
	} else {
		fmt.Println("No available prizes")
	}
}
