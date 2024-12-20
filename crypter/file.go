package crypter

import (
	"dobrify/internal/alog"
	"encoding/json"
	"log/slog"
	"os"
)

func (c *Crypter) LoadFromFile(filename string, dest interface{}) {
	body, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("file not found", "filename", filename, alog.Error(err))
		return
	}

	if len(body) == 0 {
		slog.Error("file is empty", "filename", filename)
		return
	}

	decrypted, err := c.Decrypt(body)
	if err != nil {
		slog.Error("failed to decrypt file", "filename", filename, alog.Error(err))
		return
	}
	if err := json.Unmarshal(decrypted, &dest); err != nil {
		slog.Error("failed to unmarshal decrypted file", "filename", filename, alog.Error(err))
		return
	}
}

func (c *Crypter) SaveToFile(filename string, source interface{}) error {
	marshaled, err := json.Marshal(source)
	if err != nil {
		slog.Error("failed to marshal source", "filename", filename, alog.Error(err))
		return err
	}
	encrypted, err := c.Encrypt(marshaled)
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
