package storage

import (
	"encoding/json"
	"fmt"
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
		return fmt.Errorf("failed to marshal source: %w", err)
	}
	if err := os.WriteFile(filename, marshaled, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
