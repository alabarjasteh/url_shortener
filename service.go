package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type (
	Shortener interface {
		PostURL(ctx context.Context, originalURL string) (string, error)
		GetURL(ctx context.Context, shortURL string) (string, error)
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
type shortenerService struct {
	cache Cache
	repo  Repository
}

func NewShortenerService(cache Cache, repo Repository) Shortener {
	return &shortenerService{
		cache: cache,
		repo:  repo,
	}
}

func (svc *shortenerService) PostURL(ctx context.Context, originalURL string) (string, error) {
	shortURL := GenerateShortLink(originalURL)

	// write through
	err := svc.cache.Set(shortURL, originalURL)
	if err != nil {
		return "", err
	}

	err = svc.repo.Store(shortURL, originalURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (svc *shortenerService) GetURL(ctx context.Context, shortURL string) (string, error) {
	// cache-aside
	originalURL, err := svc.cache.Get(shortURL)
	if err == redis.Nil {
		// does not exist in cache
		// retrive from DB
		originalURL, err = svc.repo.Load(shortURL)
		if err != nil {
			return "", err
		}

		err = svc.cache.Set(shortURL, originalURL)
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
