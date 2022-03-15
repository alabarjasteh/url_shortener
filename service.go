package main

import (
	"github.com/alabarjasteh/url-shortener/db"
	"github.com/alabarjasteh/url-shortener/memcache"
)

type ShortenerService struct {
	db       db.Interface
	memcache memcache.Interface
}

func NewShortenerService(db db.Interface, cache memcache.Interface) *ShortenerService {
	return &ShortenerService{db: db, memcache: cache}
}
