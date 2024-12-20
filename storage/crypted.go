package storage

import (
	"dobrify/crypter"
	"dobrify/internal/alog"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type cryptedStorage struct {
	crtpr *crypter.Crypter
}

func NewCryptedStore(key string) Storage {
	return &cryptedStorage{
		crtpr: crypter.NewCrypter(key),
	}
}

func (c *cryptedStorage) LoadFromFile(filename string, dest interface{}) error {
	body, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToReadFile, err)
	}

	if len(body) == 0 {
		return ErrFileIsEmpty
	}

	decrypted, err := c.crtpr.Decrypt(body)
	if err != nil {
		return fmt.Errorf("failed to decrypt file: %w", err)
	}
	if err := json.Unmarshal(decrypted, &dest); err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}
	return nil
}

func (c *cryptedStorage) SaveToFile(filename string, source interface{}) error {
	marshaled, err := json.Marshal(source)
	if err != nil {
		slog.Error("failed to marshal source", "filename", filename, alog.Error(err))
		return err
	}
	encrypted, err := c.crtpr.Encrypt(marshaled)
	if err != nil {
		slog.Error("failed to encrypt source", "filename", filename, alog.Error(err))
		return err
	}
	if err := os.WriteFile(filename, encrypted, 0644); err != nil {
		slog.Error("failed to write file", "filename", filename, alog.Error(err))
		return err
	}
	return nil
}
