package middleware

import (
	"booking/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	rdb *redis.Client
	log logger.Logger
}

func NewMiddlewares(rdb *redis.Client, log logger.Logger) *Middleware {
	return &Middleware{rdb: rdb, log: log}
}
