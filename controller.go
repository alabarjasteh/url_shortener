package main

import (
	"github.com/alabarjasteh/url-shortener/db"
	"github.com/alabarjasteh/url-shortener/memcache"
)

type Controller struct {
	db       db.Interface
	memcache memcache.Interface
}

func NewController(db db.Interface, cache memcache.Interface) *Controller {
	return &Controller{db: db, memcache: cache}
}
