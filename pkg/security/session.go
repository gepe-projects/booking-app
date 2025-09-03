package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"booking/internal/domain"
	"booking/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const (
	sessionTTL       = 1 * time.Minute    // TTL sliding session
	sessionExtendTTL = 30 * time.Minute   // extend kalau sisa < 30 menit
	userSessionTTL   = 7 * 24 * time.Hour // kalo user cuma coba coba login, dihapus setelah 30 hari
)

// =============================
// CREATE & EXTEND SESSION
// =============================

func CreateSession(
	ctx context.Context,
	rdb *redis.Client,
	data *domain.UserWithIdentity,
	log logger.Logger,
	device string,
	userAgent string,
	ipAddress string,
) (string, error) {
	token, err := generateOpaqueToken(32)
	if err != nil {
		log.Error(err, "failed to generate opaque token")
		return "", err
	}

	sessionKey := generateSessionKey(token)
	userSessionsKey := generateUserSessionsKey(data.User.ID)

	sessionValues := domain.Session{
		UserID:     data.User.ID,
		Name:       data.User.Name,
		Role:       data.User.Role,
		ImageURL:   domain.NilStringHandler(data.User.ImageURL),
		Email:      domain.NilStringHandler(data.UserIdentity.Email),
		Provider:   data.UserIdentity.Provider,
		ProviderID: data.UserIdentity.ProviderID,
		Phone:      domain.NilStringHandler(data.UserIdentity.Phone),
		Verified:   domain.BoolToString(data.UserIdentity.Verified),
		Device:     device,
		UserAgent:  userAgent,
		IpAddress:  ipAddress,
	}

	expireTime := time.Now().Add(sessionTTL).Unix()

	pipe := rdb.Pipeline()
	pipe.HSet(ctx, sessionKey, sessionValues.ToRedisMap())
	pipe.Expire(ctx, sessionKey, sessionTTL)
	pipe.ZAdd(ctx, userSessionsKey, redis.Z{
		Score:  float64(expireTime),
		Member: token,
	})
	pipe.Expire(ctx, userSessionsKey, userSessionTTL)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error(err, "failed to create session pipeline")
		return "", err
	}

	return token, nil
}

func ExtendSession(ctx context.Context, rdb *redis.Client, token string, userID string) error {
	sessionKey := generateSessionKey(token)
	userSessionsKey := generateUserSessionsKey(userID)

	exists, err := rdb.Exists(ctx, sessionKey).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return domain.ErrUnauthorized
	}

	newExpireTime := time.Now().Add(sessionTTL).Unix()

	pipe := rdb.Pipeline()
	pipe.Expire(ctx, sessionKey, sessionTTL)
	pipe.ZAdd(ctx, userSessionsKey, redis.Z{
		Score:  float64(newExpireTime),
		Member: token,
	})
	pipe.Expire(ctx, userSessionsKey, userSessionTTL)
	_, err = pipe.Exec(ctx)
	return err
}

// =============================
// GET SESSION
// =============================

func GetSession(ctx context.Context, rdb *redis.Client, token string, log logger.Logger) (*domain.Session, error) {
	sessionKey := generateSessionKey(token)

	pipe := rdb.Pipeline()
	hgetallCmd := pipe.HGetAll(ctx, sessionKey)
	ttlCmd := pipe.TTL(ctx, sessionKey)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Error(err, "failed to exec GetSession pipeline")
		return nil, domain.ErrInternalServerError
	}

	if len(hgetallCmd.Val()) == 0 {
		log.Error(nil, "session not found")
		return nil, domain.ErrUnauthorized
	}

	var session domain.Session
	if err := hgetallCmd.Scan(&session); err != nil {
		log.Error(err, "failed to scan hgetallCmd to session")
		return nil, domain.ErrInternalServerError
	}

	if ttlCmd.Val() < sessionExtendTTL && ttlCmd.Val() > 0 {
		if err := ExtendSession(ctx, rdb, token, session.UserID); err != nil {
			log.Error(err, "failed to extend session")
		}
	}
	return &session, nil
}

// =============================
// GET USER ACTIVE SESSIONS (CLEANUP AUTO)
// =============================

func GetUserActiveSessionsWithDetails(ctx context.Context, rdb *redis.Client, userID string, log logger.Logger) ([]domain.SessionWithExpiry, error) {
	userSessionsKey := generateUserSessionsKey(userID)
	now := time.Now().Unix()

	// Hapus expired tokens dulu
	_, _ = rdb.ZRemRangeByScore(ctx, userSessionsKey, "-inf", fmt.Sprintf("%d", now)).Result()

	// Ambil tokens dengan score (expire time)
	tokensWithScores, err := rdb.ZRangeByScoreWithScores(ctx, userSessionsKey, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", now),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, err
	}

	if len(tokensWithScores) == 0 {
		return []domain.SessionWithExpiry{}, nil
	}

	var sessionsWithExpiry []domain.SessionWithExpiry
	var zombieTokens []string

	for _, tokenScore := range tokensWithScores {
		token := tokenScore.Member.(string)
		expireTime := int64(tokenScore.Score)

		sessionKey := generateSessionKey(token)
		sessionData, err := rdb.HGetAll(ctx, sessionKey).Result()
		if err != nil || len(sessionData) == 0 {
			zombieTokens = append(zombieTokens, token)
			continue
		}

		var session domain.Session
		if err := mapToSession(sessionData, &session); err != nil {
			log.Error(err, "failed to convert session data")
			continue
		}

		sessionsWithExpiry = append(sessionsWithExpiry, domain.SessionWithExpiry{
			Session:    session,
			ExpireTime: time.Unix(expireTime, 0),
			Token:      token,
		})
	}

	// Hapus zombie tokens dari ZSET
	if len(zombieTokens) > 0 {
		_, _ = rdb.ZRem(ctx, userSessionsKey, zombieTokens).Result()
		log.Infof("cleaned %d zombie tokens for user:%s", len(zombieTokens), userID)
	}

	return sessionsWithExpiry, nil
}

// =============================
// LOGOUT
// =============================

func LogoutSession(ctx context.Context, rdb *redis.Client, userID, token string) error {
	sessionKey := generateSessionKey(token)
	userSessionsKey := generateUserSessionsKey(userID)

	pipe := rdb.Pipeline()
	pipe.Del(ctx, sessionKey)
	pipe.ZRem(ctx, userSessionsKey, token)
	_, err := pipe.Exec(ctx)
	return err
}

func LogoutAllSessions(ctx context.Context, rdb *redis.Client, userID string) error {
	userSessionsKey := generateUserSessionsKey(userID)
	tokens, err := rdb.ZRange(ctx, userSessionsKey, 0, -1).Result()
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return nil
	}

	pipe := rdb.Pipeline()
	for _, token := range tokens {
		pipe.Del(ctx, generateSessionKey(token))
	}
	pipe.Del(ctx, userSessionsKey)
	_, err = pipe.Exec(ctx)
	return err
}

func LogoutOtherSessions(ctx context.Context, rdb *redis.Client, userID, currentToken string) error {
	userSessionsKey := generateUserSessionsKey(userID)
	tokens, err := rdb.ZRange(ctx, userSessionsKey, 0, -1).Result()
	if err != nil {
		return err
	}

	pipe := rdb.Pipeline()
	for _, token := range tokens {
		if token == currentToken {
			continue
		}
		pipe.Del(ctx, generateSessionKey(token))
		pipe.ZRem(ctx, userSessionsKey, token)
	}
	_, err = pipe.Exec(ctx)
	return err
}

// =============================
// HELPERS
// =============================

func mapToSession(data map[string]string, session *domain.Session) error {
	session.UserID = data["userID"]
	session.Name = data["name"]
	session.Role = data["role"]
	session.ImageURL = data["image_url"]
	session.Email = data["email"]
	session.Provider = data["provider"]
	session.ProviderID = data["provider_id"]
	session.Phone = data["phone"]
	session.Verified = data["verified"]
	session.Device = data["device"]
	session.UserAgent = data["user_agent"]
	session.IpAddress = data["ip_address"]
	return nil
}

func generateSessionKey(token string) string {
	return fmt.Sprintf("session:%s", token)
}

func generateUserSessionsKey(userID string) string {
	return fmt.Sprintf("user_sessions:%s", userID)
}

func generateOpaqueToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
