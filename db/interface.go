package db

import "errors"

type Interface interface {
	Load(shortLink string) (string, error)
	Save(shortLink, originalLink string) error
}

var (
	ErrNotFound = errors.New("not found")
)
