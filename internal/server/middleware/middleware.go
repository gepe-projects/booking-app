package middleware

import (
	"booking/pkg/logger"
	"booking/pkg/security"

	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	security *security.Security
	rdb      *redis.Client
	log      logger.Logger
}

func NewMiddlewares(security *security.Security, rdb *redis.Client, log logger.Logger) *Middleware {
	return &Middleware{security: security, rdb: rdb, log: log}
}
