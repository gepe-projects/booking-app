package security

import (
	"booking/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type Security struct {
	rdb *redis.Client
	log logger.Logger
}

func NewSecurity(rdb *redis.Client, log logger.Logger) *Security {
	return &Security{rdb: rdb, log: log}
}
