package memcache

import (
	"context"
	"fmt"

	"github.com/alabarjasteh/url-shortener/config"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

var ctx = context.Background() //TODO

func NewRedis(c *config.Config) Interface {
	addr := fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
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
		// panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
		return err
	}
	return nil
}

func (r *Redis) Get(shortUrl string) (string, error) {
	result, err := r.Client.Get(ctx, shortUrl).Result()
	if err != nil {
		// panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, shortUrl))
		return "", err
	}
	return result, nil
}
