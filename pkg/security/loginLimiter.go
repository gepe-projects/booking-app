package security

import (
	"context"
	"fmt"
	"time"

	"booking/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const (
	loginAttemptsKey = "login_attempts"
	loginBanKey      = "login_ban"
)

// CheckBan: apakah user sedang diban?
func CheckBan(ctx context.Context, rdb *redis.Client, email string, log logger.Logger) (time.Duration, error) {
	banKey := fmt.Sprintf("%s:%s", loginBanKey, email)

	ttl, err := rdb.TTL(ctx, banKey).Result()
	if err != nil && err != redis.Nil {
		log.Error(err, "failed to get ban ttl in redis")
		return 0, err
	}
	if ttl > 0 {
		return ttl, nil
	}
	return 0, nil
}

// IncrementAttempts: dipanggil kalau login gagal
func IncrementAttempts(ctx context.Context, rdb *redis.Client, email string, log logger.Logger) (time.Duration, error) {
	attemptsKey := fmt.Sprintf("%s:%s", loginAttemptsKey, email)
	banKey := fmt.Sprintf("%s:%s", loginBanKey, email)

	// increment attempts
	attempts, err := rdb.Incr(ctx, attemptsKey).Result()
	if err != nil {
		log.Error(err, "failed to increment login attempts in redis")
		return 0, err
	}

	// TTL panjang buat counter (misalnya 1 jam sejak first fail)
	if attempts == 1 {
		_ = rdb.Expire(ctx, attemptsKey, 1*time.Hour).Err()
	}

	// kalau attempts > 3 → hitung delay ban
	if attempts > 3 {
		delay := time.Duration(5*(1<<(attempts-3))) * time.Second
		if delay > 1*time.Hour {
			delay = 1 * time.Hour
		}

		// set ban key dengan TTL delay
		if err := rdb.SetEx(ctx, banKey, "1", delay).Err(); err != nil {
			log.Error(err, "failed to set ban key in redis")
			return 0, err
		}
		return delay, nil
	}

	return 0, nil
}

// ResetLoginAttempts: dipanggil kalau login sukses
func ResetLoginAttempts(ctx context.Context, rdb *redis.Client, email string, log logger.Logger) error {
	attemptsKey := fmt.Sprintf("%s:%s", loginAttemptsKey, email)
	banKey := fmt.Sprintf("%s:%s", loginBanKey, email)

	_, err := rdb.Del(ctx, attemptsKey, banKey).Result()
	if err != nil {
		log.Error(err, "failed to reset login attempts in redis")
		return err
	}
	return nil
}
