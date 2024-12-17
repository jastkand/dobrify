package crypter

import (
	"encoding/json"
	"log/slog"
	"os"
)

func LoadFromFile(secretKey, filename string, dest interface{}) {
	body, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("file not found", "filename", filename, "error", err.Error())
		return
	}

	if len(body) == 0 {
		slog.Error("file is empty", "filename", filename)
		return
	}

	cpt := NewCrypter(secretKey)
	decrypted, err := cpt.Decrypt(body)
	if err != nil {
		slog.Error("failed to decrypt file", "filename", filename, "error", err.Error())
		return
	}
	if err := json.Unmarshal(decrypted, &dest); err != nil {
		slog.Error("failed to unmarshal decrypted file", "filename", filename, "error", err.Error())
		return
	}
}

func SaveToFile(secretKey, filename string, source interface{}) error {
	cpt := NewCrypter(secretKey)
	marshaled, err := json.Marshal(source)
	if err != nil {
		slog.Error("failed to marshal source", "filename", filename, "error", err.Error())
		return err
	}
	encrypted, err := cpt.Encrypt(marshaled)
	if err != nil {
		slog.Error("failed to encrypt source", "filename", filename, "error", err.Error())
		return err
	}
	if err := os.WriteFile(filename, encrypted, 0644); err != nil {
		slog.Error("failed to write file", "filename", filename, "error", err.Error())
		return err
	}
	return nil
}
