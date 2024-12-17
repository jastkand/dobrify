package dobry

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const mimeJSON = "application/json"
const origin = "https://dobrycola-promo.ru"
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'"

const BaseURL = "https://dobrycola-promo.ru/backend"

type Client struct {
	HttpClient *http.Client
	Username   string
	Password   string
	Token      *Token
}

func NewClient(username, password string, token *Token) *Client {
	return &Client{
		HttpClient: http.DefaultClient,
		Username:   username,
		Password:   password,
		Token:      token,
	}
}

type ensureTokenResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Token  `json:"data"`
}

func (c *Client) EnsureToken() (*Token, error) {
	if c.Token != nil {
		slog.Debug("checking existing token")
		if c.isAccessTokenValid() {
			return c.Token, nil
		}
		slog.Debug("token expired, refreshing")
		err := c.RefreshToken()
		if err == nil && c.isAccessTokenValid() {
			return c.Token, nil
		}
	}

	slog.Debug("getting new token")

	req, err := http.NewRequest(http.MethodPost, BaseURL+"/oauth/token", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", mimeJSON)
	req.Header.Set("Content-Type", mimeJSON)
	req.Header.Set("Origin", origin)
	req.Header.Set("Referer", "https://dobrycola-promo.ru/?signin")
	req.Header.Set("User-Agent", userAgent)

	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, c.Username, c.Password)))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get token: %d", resp.StatusCode)
	}

	var parsed ensureTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	c.Token = &parsed.Data
	return c.Token, nil
}

type refreshResponse struct {
	Status int   `json:"status"`
	Data   Token `json:"data"`
}

func (c *Client) RefreshToken() error {
	req, err := http.NewRequest(http.MethodPost, BaseURL+"/oauth/refresh", nil)
	if err != nil {
		slog.Error("failed to create refresh token request", "error", err.Error())
		return err
	}
	req.Header.Set("Accept", mimeJSON)
	req.Header.Set("Content-Type", mimeJSON)
	req.Header.Set("Origin", origin)
	req.Header.Set("Referer", "https://dobrycola-promo.ru/profile")
	req.Header.Set("User-Agent", userAgent)

	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"refresh_token":"%s"}`, c.Token.RefreshToken)))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		slog.Error("failed to send refresh token request", "error", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to refresh token", "status", resp.StatusCode)
		return fmt.Errorf("failed to refresh token: %d", resp.StatusCode)
	}

	var parsed refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		slog.Error("failed to parse refresh token response", "error", err.Error())
		return err
	}

	c.Token = &parsed.Data
	return nil
}

type ShopResponse struct {
	Status int        `json:"status"`
	Data   []ShopItem `json:"data"`
}

func (c *Client) GetPrizes() (*ShopResponse, error) {
	req, err := http.NewRequest(http.MethodGet, BaseURL+"/private/prize/shop", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", mimeJSON)
	req.Header.Set("Authorization", "Bearer "+c.Token.AccessToken)
	req.Header.Set("Referer", "https://dobrycola-promo.ru/profile")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get prizes: %d", resp.StatusCode)
	}

	var parsed ShopResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

func (c *Client) GetProfile() error {
	req, err := http.NewRequest(http.MethodGet, BaseURL+"/private/user", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", mimeJSON)
	req.Header.Set("Authorization", "Bearer "+c.Token.AccessToken)
	req.Header.Set("Referer", "https://dobrycola-promo.ru/?signin")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get profile: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) isAccessTokenValid() bool {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(c.Token.AccessToken, jwt.MapClaims{})
	if err != nil {
		slog.Error("auth token is invalid", "error", err.Error())
		return false
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		slog.Error("failed to get expiration time", "error", err.Error())
		return false
	}

	if time.Now().After(exp.Time) {
		slog.Error("token is expired", "exp", exp.Format(time.RFC3339))
		return false
	}

	return true
}
