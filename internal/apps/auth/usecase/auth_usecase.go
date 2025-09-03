package usecase

import (
	"context"

	"booking/internal/domain"
	"booking/pkg/logger"
	"booking/pkg/security"

	"github.com/redis/go-redis/v9"
)

type authUsecase struct {
	userUsecase domain.UserUsecase
	rdb         *redis.Client
	log         logger.Logger
}

func NewAuthUsecase(userUsecase domain.UserUsecase, rdb *redis.Client, log logger.Logger) domain.AuthUsecase {
	return &authUsecase{
		userUsecase: userUsecase,
		rdb:         rdb,
		log:         log,
	}
}

func (u *authUsecase) RegisterUser(ctx context.Context, req *domain.RegisterDTO) (*domain.UserWithIdentity, error) {
	res, err := u.userUsecase.RegisterUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginDTO) (*domain.UserWithIdentity, string, error) {
	// get user by email
	res, token, err := u.userUsecase.Login(ctx, req)
	if err != nil {
		return nil, "", err
	}

	return res, token, nil
}

func (u *authUsecase) GetAllActiveSessions(ctx context.Context, userId string) ([]domain.SessionWithExpiry, error) {
	return security.GetUserActiveSessionsWithDetails(ctx, u.rdb, userId, u.log)
}
