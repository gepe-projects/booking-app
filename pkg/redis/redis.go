package redis

import (
	"context"
	"fmt"
	"time"

	"booking/pkg/config"
	"booking/pkg/logger"

	"github.com/redis/go-redis/v9"
)

func NewClient(redisCfg *config.RedisConfig, log logger.Logger) *redis.Client {
	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port),
		Password:     redisCfg.Password,
		DB:           0,
		PoolSize:     100, // tune according to concurrency
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  30 * time.Second,
	}

	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err, "error ping redis")
	}
	log.Info("redis connected")
	return rdb
}
