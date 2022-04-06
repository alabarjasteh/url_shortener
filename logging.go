package main

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

type Middleware func(Shortener) Shortener

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Shortener) Shortener {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Shortener
	logger log.Logger
}

func (mw loggingMiddleware) PostURL(ctx context.Context, repo Repository, cache Cache, originalURL string) (string, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostURL", "took", time.Since(begin))
	}(time.Now())
	return mw.next.PostURL(ctx, repo, cache, originalURL)
}

func (mw loggingMiddleware) GetURL(ctx context.Context, repo Repository, cache Cache, shortURL string) (string, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetURL", "took", time.Since(begin))
	}(time.Now())
	return mw.next.GetURL(ctx, repo, cache, shortURL)
}
