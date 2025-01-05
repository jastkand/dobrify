package crypter

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCrypterEncrypt(t *testing.T) {
	t.Parallel()
	t.Run("encrypts text", func(t *testing.T) {
		t.Parallel()
		c := NewCrypter(strings.ReplaceAll(uuid.NewString(), "-", ""))
		encrypted, err := c.Encrypt([]byte("text"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if string(encrypted) == "text" {
			t.Errorf("expected text to be encrypted")
		}
	})
}

func TestCrypterDecrypt(t *testing.T) {
	t.Parallel()
	t.Run("decrypts text", func(t *testing.T) {
		t.Parallel()
		c := NewCrypter(strings.ReplaceAll(uuid.NewString(), "-", ""))
		encrypted, _ := c.Encrypt([]byte("text"))
		decrypted, err := c.Decrypt(encrypted)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if string(decrypted) != "text" {
			t.Errorf("expected text to be decrypted")
		}
	})
}
