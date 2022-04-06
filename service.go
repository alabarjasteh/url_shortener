package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type (
	Shortener interface {
		PostURL(ctx context.Context, repo Repository, cache Cache, originalURL string) (string, error)
		GetURL(ctx context.Context, repo Repository, cache Cache, shortURL string) (string, error)
	}
	Repository interface {
		Load(shortLink string) (string, error)
		Store(shortLink, originalLink string) error
	}
	Cache interface {
		Get(shortLink string) (string, error)
		Set(shortLink, originalLink string) error
	}
)

// concrete object that implements Shortener interface
type shortenerService struct{}

func NewShortenerService() Shortener {
	return &shortenerService{}
}

func (svc *shortenerService) PostURL(ctx context.Context, repo Repository, cache Cache, originalURL string) (string, error) {
	shortURL := GenerateShortLink(originalURL)

	// write through
	err := cache.Set(shortURL, originalURL)
	if err != nil {
		return "", err
	}

	err = repo.Store(shortURL, originalURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (svc *shortenerService) GetURL(ctx context.Context, repo Repository, cache Cache, shortURL string) (string, error) {
	// cache-aside
	originalURL, err := cache.Get(shortURL)
	if err == redis.Nil {
		// does not exist in cache
		// retrive from DB
		originalURL, err = repo.Load(shortURL)
		if err != nil {
			return "", err
		}

		err = cache.Set(shortURL, originalURL)
		if err != nil {
			return "", err
		}
		return originalURL, nil
	}
	if err != nil {
		return "", err
	}
	return originalURL, nil
}
