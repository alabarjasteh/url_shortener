package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alabarjasteh/url-shortener/config"
	"github.com/alabarjasteh/url-shortener/db"
	"github.com/alabarjasteh/url-shortener/memcache"
	"github.com/alabarjasteh/url-shortener/shortener"
	"github.com/go-kit/log"
)

func main() {
	configPath := "./config/config"
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(file)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	db := db.NewMySql(cfg)
	cache := memcache.NewRedis(cfg)

	var svc shortener.Interface
	svc = NewShortenerService(db, cache)
	svc = LoggingMiddleware(logger)(svc)
	router := MakeRoutes(svc)

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), router)
}
