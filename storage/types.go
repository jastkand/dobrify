package storage

import "errors"

var (
	ErrFailedToReadFile = errors.New("failed to read file")
	ErrFileIsEmpty      = errors.New("file is empty")
)

type Storage interface {
	LoadFromFile(filename string, dest interface{}) error
	SaveToFile(filename string, source interface{}) error
}
