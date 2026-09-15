package db

import (
	"Orbit/configs"
	"log"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func SetUpRedis() {
	cfg := configs.LoadConfig().Redis
	options := &redis.Options{Addr: cfg.Url, DB: 0}

	if strings.Contains(cfg.Url, "://") {
		parsed, err := redis.ParseURL(cfg.Url)
		if err != nil {
			log.Fatalf("invalid REDIS_URL: %v", err)
		}
		options = parsed
	}

	client := redis.NewClient(options)

	redisClient = client
}

func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		SetUpRedis()
	})
	return redisClient
}
