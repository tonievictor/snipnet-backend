package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Init(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return rdb
}
