package usecase

import (
	"context"
	"database/sql"
	"errors"

	"booking/internal/domain"
	"booking/pkg/logger"
)

type userUseCase struct {
	userRepository domain.UserRepository
	log            logger.Logger
}

func NewUserUseCase(userRepository domain.UserRepository, log logger.Logger) domain.UserUsecase {
	return &userUseCase{
		userRepository: userRepository,
		log:            log,
	}
}

func (u *userUseCase) GetByEmail(ctx context.Context, email string) (*domain.UserWithIdentity, error) {
	res, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.log.Debug("User not found")
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrInternalServerError
	}

	return res, nil
}
