package security

import (
	"booking/pkg/config"
	"booking/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type Security struct {
	config *config.Config
	rdb    *redis.Client
	log    logger.Logger
}

func NewSecurity(config *config.Config, rdb *redis.Client, log logger.Logger) *Security {
	return &Security{config: config, rdb: rdb, log: log}
}
