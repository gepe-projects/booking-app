package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	loginAttemptsKey = "login_attempts"
	loginBanKey      = "login_ban"
)

// CheckBan: apakah user sedang diban?
func (s *Security) CheckBan(ctx context.Context, email string) (time.Duration, error) {
	banKey := fmt.Sprintf("%s:%s", loginBanKey, email)

	ttl, err := s.rdb.TTL(ctx, banKey).Result()
	if err != nil && err != redis.Nil {
		s.log.Error(err, "failed to get ban ttl in redis")
		return 0, err
	}
	if ttl > 0 {
		return ttl, nil
	}
	return 0, nil
}

// IncrementAttempts: dipanggil kalau login gagal
func (s *Security) IncrementAttempts(ctx context.Context, email string) (time.Duration, error) {
	attemptsKey := fmt.Sprintf("%s:%s", loginAttemptsKey, email)
	banKey := fmt.Sprintf("%s:%s", loginBanKey, email)

	// increment attempts
	attempts, err := s.rdb.Incr(ctx, attemptsKey).Result()
	if err != nil {
		s.log.Error(err, "failed to increment login attempts in redis")
		return 0, err
	}

	// TTL panjang buat counter (misalnya 1 jam sejak first fail)
	if attempts == 1 {
		_ = s.rdb.Expire(ctx, attemptsKey, 1*time.Hour).Err()
	}

	// kalau attempts > 3 → hitung delay ban
	if attempts > 3 {
		delay := time.Duration(5*(1<<(attempts-3))) * time.Second
		if delay > 1*time.Hour {
			delay = 1 * time.Hour
		}

		// set ban key dengan TTL delay
		if err := s.rdb.SetEx(ctx, banKey, "1", delay).Err(); err != nil {
			s.log.Error(err, "failed to set ban key in redis")
			return 0, err
		}
		return delay, nil
	}

	return 0, nil
}

// ResetLoginAttempts: dipanggil kalau login sukses
func (s *Security) ResetLoginAttempts(ctx context.Context, email string) error {
	attemptsKey := fmt.Sprintf("%s:%s", loginAttemptsKey, email)
	banKey := fmt.Sprintf("%s:%s", loginBanKey, email)

	_, err := s.rdb.Del(ctx, attemptsKey, banKey).Result()
	if err != nil {
		s.log.Error(err, "failed to reset login attempts in redis")
		return err
	}
	return nil
}
