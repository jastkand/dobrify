package dobry

import (
	"dobrify/crypter"
	"encoding/json"
	"log/slog"
	"os"
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
)

type App struct {
	Client    *Client
	SecretKey string
}

func NewApp(username, password, secretKey string) *App {
	token := loadTokenFromFile(secretKey)
	return &App{
		Client:    NewClient(username, password, token),
		SecretKey: secretKey,
	}
}

func (a *App) HasWantedPrizes(wantedPrizes []string) ([]string, error) {
	token, err := a.Client.EnsureToken()
	if err != nil {
		slog.Error("failed to ensure token", "error", err.Error())
		return nil, err
	}
	if err := saveTokenToFile(token, a.SecretKey); err != nil {
		return nil, err
	}
	prizes, err := a.Client.GetPrizes()
	if err != nil {
		slog.Error("failed to get prizes", "error", err.Error())
		return nil, err
	}
	slog.Info("got prizes", "prizes_count", len(prizes.Data))
	var availablePrizes []string
	for _, prize := range prizes.Data {
		if isWantedPrize(wantedPrizes, prize.Alias) && !prize.TotalLimit {
			availablePrizes = append(availablePrizes, prize.Alias)
		}
	}
	return availablePrizes, nil
}

func isWantedPrize(wantedPrizes []string, prize string) bool {
	for _, wanted := range wantedPrizes {
		if prize == wanted {
			return true
		}
	}
	return false
}

func loadTokenFromFile(secretKey string) *Token {
	body, err := os.ReadFile("tokens.bin")
	if err != nil {
		slog.Error("no tokens.bin file was found", "error", err.Error())
		return nil
	}

	if len(body) == 0 {
		slog.Error("tokens.bin file is empty")
		return nil
	}

	var token *Token
	cpt := crypter.NewCrypter(secretKey)
	decrypted, err := cpt.Decrypt(body)
	if err != nil {
		slog.Error("failed to decrypt tokens.bin", "error", err.Error())
		return nil
	}
	if err := json.Unmarshal(decrypted, &token); err != nil {
		slog.Error("failed to unmarshal decrypted tokens.bin", "error", err.Error())
		return nil
	}
	return token
}

func saveTokenToFile(token *Token, secretKey string) error {
	cpt := crypter.NewCrypter(secretKey)
	marshaled, err := json.Marshal(token)
	if err != nil {
		slog.Error("failed to marshal token", "error", err.Error())
		return err
	}
	encrypted, err := cpt.Encrypt(marshaled)
	if err != nil {
		slog.Error("failed to encrypt token", "error", err.Error())
		return err
	}
	if err := os.WriteFile("tokens.bin", encrypted, 0644); err != nil {
		slog.Error("failed to write token to file", "error", err.Error())
		return err
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
