package storage

import (
	"encoding/json"
	"fmt"
)

type inMemoryStorage struct {
	data []byte
}

func NewInMemoryStore(data []byte) Storage {
	return &inMemoryStorage{
		data: data,
	}
}

func (s *inMemoryStorage) LoadFromFile(_ string, dest interface{}) error {
	if len(s.data) == 0 {
		s.data = []byte("{}")
	}

	if err := json.Unmarshal(s.data, &dest); err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}
	return nil
}

func (s *inMemoryStorage) SaveToFile(_ string, source interface{}) error {
	marshaled, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal source: %w", err)
	}
	s.data = marshaled
	return nil
}
