package rediscache

import (
	"context"
	"fmt"

	"github.com/alabarjasteh/url-shortener/config"
	urlshortenersvc "github.com/alabarjasteh/url-shortener/urlshortener"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

var ctx = context.Background() //TODO

func New(c *config.RedisConfig) *Redis {
	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: c.Password,
		DB:       c.DB,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}
	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)

	return &Redis{
		Client: redisClient,
	}
}

func (r *Redis) Set(shortUrl, originalUrl string) error {
	err := r.Client.Set(ctx, shortUrl, originalUrl, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(shortUrl string) (string, error) {
	result, err := r.Client.Get(ctx, shortUrl).Result()
	if err != nil {
		if err == redis.Nil {
			return "", urlshortenersvc.ErrCacheMiss
		}
		return "", err
	}
	return result, nil
}
