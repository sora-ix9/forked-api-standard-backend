package redisclient

import (
	"fdlp-standard-api/pkg/config"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func NewClient(cfg *config.Config) *redis.Client {
	var db int = 0
	db, _ = strconv.Atoi(cfg.RedisDb)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort, // Redis server address
		Password: cfg.RedisPassword,                   // no password set
		DB:       db,                                  // use default DB
	})

	return rdb
}
