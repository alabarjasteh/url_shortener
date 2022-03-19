package main

import (
	"net/http"
	"time"

	"github.com/alabarjasteh/url-shortener/shortener"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

type Middleware func(shortener.Interface) shortener.Interface

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next shortener.Interface) shortener.Interface {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   shortener.Interface
	logger log.Logger
}

func (mw loggingMiddleware) PostUrl(w http.ResponseWriter, req *http.Request) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostUrl", "took", time.Since(begin))
	}(time.Now())
	mw.next.PostUrl(w, req)
}

func (mw loggingMiddleware) RedirectShortUrl(w http.ResponseWriter, req *http.Request) {
	defer func(begin time.Time) {
		params := mux.Vars(req)
		shortUrl := params["shortlink"]
		mw.logger.Log("method", "RedirectShortUrl", "shortLink", shortUrl, "took", time.Since(begin))
	}(time.Now())
	mw.next.RedirectShortUrl(w, req)
}
