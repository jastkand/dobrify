package dobry

import (
	"dobrify/internal/alog"
	"dobrify/internal/config"
	"dobrify/storage"
	"fmt"
	"log/slog"
)

const (
	GlassRed        = "shop_glass_red"
	GlassGreen      = "shop_glass_green"
	CapGreen        = "shop_cap_green"
	CapWhite        = "shop_cap_white"
	CapRed          = "shop_cap_red"
	ElkRed          = "shop_elk_red"
	ElkWhite        = "shop_elk_white"
	ElkGreen        = "shop_elk_green"
	SweatshirtWhite = "shop_sweatshirt_white"
	SweatshirtGreen = "shop_sweatshirt_green"
	SweatshirtRed   = "shop_sweatshirt_red"
)

var (
	Glasses     = []string{GlassRed, GlassGreen}
	Caps        = []string{CapGreen, CapWhite, CapRed}
	Elks        = []string{ElkRed, ElkWhite, ElkGreen}
	Sweatshirts = []string{SweatshirtWhite, SweatshirtGreen, SweatshirtRed}
	AllPrizes   = append(Glasses, append(Caps, append(Elks, Sweatshirts...)...)...)
)

type App struct {
	cfg      config.Config
	store    storage.Storage
	encStore storage.Storage
	dobryApi *Client
}

func NewApp(cfg config.Config) *App {
	encStore := storage.NewCryptedStore(cfg.SecretKey)

	var token *Token
	if err := encStore.LoadFromFile("tokens.bin", &token); err != nil {
		slog.Error("failed to load token", alog.Error(err))
	}

	return &App{
		cfg:      cfg,
		encStore: encStore,
		dobryApi: NewClient(cfg.DobryUsername, cfg.DobryPassword, token),
	}
}

func (a *App) GetAvailablePrizes() ([]string, error) {
	slog.Debug("getting available prizes")
	if err := a.renewToken(); err != nil {
		slog.Error("failed to renew token", alog.Error(err))
		return nil, err
	}
	prizes, err := a.dobryApi.GetPrizes()
	if err != nil {
		slog.Error("failed to get prizes", alog.Error(err))
		return nil, err
	}
	var availablePrizes []string
	for _, prize := range prizes.Data {
		if !prize.TotalLimit {
			availablePrizes = append(availablePrizes, prize.Alias)
		}
	}
	slog.Debug("got prizes",
		slog.Int("all_prizes_count", len(prizes.Data)),
		slog.Int("available_prizes_count", len(availablePrizes)),
		slog.Any("available_prizes", availablePrizes),
	)
	return availablePrizes, nil
}

func (a *App) renewToken() error {
	token, err := a.dobryApi.EnsureToken()
	if err != nil {
		return fmt.Errorf("failed to ensure token: %w", err)
	}
	if err := a.encStore.SaveToFile("tokens.bin", token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}

func PrizeName(prize string) string {
	switch prize {
	case GlassRed:
		return "Стеклянный стакан с красным логотипом"
	case GlassGreen:
		return "Стеклянный стакан с зеленым логотипом"
	case CapGreen:
		return "Шапка зеленая"
	case CapWhite:
		return "Шапка белая"
	case CapRed:
		return "Шапка красная"
	case ElkRed:
		return "Мягкая игрушка «лосик» красная"
	case ElkWhite:
		return "Мягкая игрушка «лосик» белая"
	case ElkGreen:
		return "Мягкая игрушка «лосик» зеленая"
	case SweatshirtWhite:
		return "Свитшот белый"
	case SweatshirtGreen:
		return "Свитшот зеленый"
	case SweatshirtRed:
		return "Свитшот красный"
	default:
		return prize
	}
}
