package storage

import (
	"dobrify/internal/alog"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type plainStorage struct {
}

func NewPlainStore() Storage {
	return &plainStorage{}
}

func (c *plainStorage) LoadFromFile(filename string, dest interface{}) error {
	body, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToReadFile, err)
	}

	if len(body) == 0 {
		return ErrFileIsEmpty
	}

	if err := json.Unmarshal(body, &dest); err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}
	return nil
}

func (c *plainStorage) SaveToFile(filename string, source interface{}) error {
	marshaled, err := json.Marshal(source)
	if err != nil {
		slog.Error("failed to marshal source", "filename", filename, alog.Error(err))
		return err
	}
	if err := os.WriteFile(filename, marshaled, 0644); err != nil {
		slog.Error("failed to write file", "filename", filename, alog.Error(err))
		return err
	}
	return nil
}
