package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alabarjasteh/url-shortener/config"
	"github.com/alabarjasteh/url-shortener/mysqldb"
	"github.com/alabarjasteh/url-shortener/rediscache"
	urlshortenersvc "github.com/alabarjasteh/url-shortener/urlshortener"
	"github.com/go-kit/log"
)

func main() {
	configPath := "./config/config"
	cfgFile, err := config.Load(configPath)
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
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	mysql := mysqldb.New(&cfg.Mysql)
	redis := rediscache.New(&cfg.Redis)

	var svc urlshortenersvc.Shortener
	{
		svc = urlshortenersvc.NewShortenerService(redis, mysql)
		svc = urlshortenersvc.LoggingMiddleware(logger)(svc)
	}

	var h http.Handler
	{
		h = urlshortenersvc.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"), mysql, redis)
	}

	http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), h)
}
