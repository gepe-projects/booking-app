package usecase

import (
	"context"
	"database/sql"
	"errors"

	"booking/internal/domain"
	"booking/pkg/constant"
	"booking/pkg/logger"
	"booking/pkg/security"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	security       *security.Security
	userRepository domain.UserRepository
	log            logger.Logger
}

func NewUserUseCase(userRepository domain.UserRepository, security *security.Security, log logger.Logger) domain.UserUsecase {
	return &userUseCase{
		userRepository: userRepository,
		security:       security,
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

func (u *userUseCase) RegisterUser(ctx context.Context, req *domain.RegisterDTO) (*domain.UserWithIdentity, error) {
	userId, err := uuid.NewV7()
	if err != nil {
		u.log.Error(err, "failed to generate uuidv7 for user")
		return nil, domain.ErrInternalServerError
	}

	identityId, err := uuid.NewV7()
	if err != nil {
		u.log.Error(err, "failed to generate uuidv7 for user identity")
		return nil, domain.ErrInternalServerError
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error(err, "failed to hash password")
		return nil, domain.ErrInternalServerError
	}

	r := &domain.RegisterDTO{
		ID:         userId.String(),
		Name:       req.Name,
		ImageURL:   "",
		Role:       "user",
		IdIdentity: identityId.String(),
		Email:      req.Email,
		Password:   string(hashedPassword),
	}

	res, err := u.userRepository.RegisterUser(ctx, r)
	if err != nil {
		u.log.Error(err, "error creating user")
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case constant.PgErrUniqueViolation:
				return nil, domain.ErrUserAlreadyExists
			}
		}
		return nil, err
	}

	return res, nil
}

func (u *userUseCase) Login(ctx context.Context, req *domain.LoginDTO) (*domain.UserWithIdentity, string, error) {
	// get user by email
	res, err := u.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		u.log.Error(err, "failed to get user by email")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", domain.ErrInvalidCredentials
		}
		return nil, "", err
	}

	// compare password
	if res.UserIdentity.PasswordHash == nil {
		return nil, "", domain.ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(*res.UserIdentity.PasswordHash), []byte(req.Password))
	if err != nil {
		_, _ = u.security.IncrementAttempts(ctx, req.Email)
		return nil, "", domain.ErrInvalidCredentials
	}
	if err := u.security.ResetLoginAttempts(ctx, req.Email); err != nil {
		return nil, "", domain.ErrInternalServerError
	}

	// create session
	token, err := u.security.CreateSession(ctx, res, req.Device, req.UserAgent, req.IpAddress)
	if err != nil {
		u.log.Error(err, "failed to create session")
		return nil, "", domain.ErrInternalServerError
	}

	return res, token, nil
}
