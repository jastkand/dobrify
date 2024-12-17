package dobry

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestClient_isAccessTokenValid(t *testing.T) {
	generateToken := func(claims jwt.Claims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		rsaKey, _ := rsa.GenerateKey(rand.Reader, 1024)
		signed, _ := token.SignedString(rsaKey)
		return signed
	}

	t.Run("valid token", func(t *testing.T) {
		token := generateToken(jwt.MapClaims{"exp": float64(time.Now().Unix() + 1000), "scope": "ROLE_USER", "sub": 13393})
		client := NewClient("username", "password", &Token{
			AccessToken: token,
		})
		if !client.isAccessTokenValid() {
			t.Error("expected token to be valid")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token := generateToken(jwt.MapClaims{"exp": float64(time.Now().Unix() - 100), "scope": "ROLE_USER", "sub": 13393})
		client := NewClient("username", "password", &Token{
			AccessToken: token,
		})
		if client.isAccessTokenValid() {
			t.Error("expected token to be invalid")
		}
	})
}
