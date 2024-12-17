package dobry

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ShopItem struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Alias      string `json:"alias"`
	Ball       int    `json:"ball"`
	Active     bool   `json:"active"`
	Limit      bool   `json:"limit"`
	TotalLimit bool   `json:"totalLimit"`
}
