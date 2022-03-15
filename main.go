package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alabarjasteh/url-shortener/config"
	"github.com/alabarjasteh/url-shortener/db"
	"github.com/alabarjasteh/url-shortener/memcache"
)

func main() {
	configPath := "./config/config"
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	db := db.NewMySql(cfg)
	cache := memcache.NewRedis(cfg)

	controller := NewController(db, cache)
	router := controller.MakeRoutes()

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), router)
}
