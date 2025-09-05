package middleware

import (
	"booking/pkg/config"
	"booking/pkg/logger"
	"booking/pkg/security"

	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	security *security.Security
	rdb      *redis.Client
	config   *config.Config
	log      logger.Logger
}

func NewMiddlewares(security *security.Security, rdb *redis.Client, config *config.Config, log logger.Logger) *Middleware {
	return &Middleware{security: security, rdb: rdb, config: config, log: log}
}
