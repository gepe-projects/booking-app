package repository

import (
	"context"
	"errors"

	"booking/internal/domain"
	"booking/pkg/constant"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.UserWithIdentity, error) {
	var res domain.UserWithIdentity

	query := `
		SELECT
			u.id AS "user.id",
			u.name AS "user.name",
			u.image_url AS "user.image_url",
			u.role AS "user.role",
			u.created_at AS "user.created_at",
			u.updated_at AS "user.updated_at",

			ui.id AS "useridentity.id",
			ui.user_id AS "useridentity.user_id",
			ui.provider AS "useridentity.provider",
			ui.provider_id AS "useridentity.provider_id",
			ui.email AS "useridentity.email",
			ui.phone AS "useridentity.phone",
			ui.password_hash AS "useridentity.password_hash",
			ui.verified AS "useridentity.verified",
			ui.created_at AS "useridentity.created_at",
			ui.updated_at AS "useridentity.updated_at"
		FROM user_identities ui
		JOIN users u ON ui.user_id = u.id
		WHERE ui.email = $1
	`

	err := r.DB.GetContext(ctx, &res, query, email)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *userRepository) RegisterUser(ctx context.Context, req *domain.RegisterDTO) (*domain.UserWithIdentity, error) {
	// start trx
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// insert ke users
	var user domain.User
	createUserQuery := `
		INSERT INTO users (id, name, image_url, role)
			VALUES ($1, $2, $3, $4)
		RETURNING id, name, image_url, role, created_at, updated_at
	`
	err = tx.QueryRowxContext(ctx, createUserQuery,
		req.ID,
		req.Name,
		req.ImageURL,
		req.Role,
	).StructScan(&user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case constant.PgErrUniqueViolation:
				return nil, err
			}
		}
		return nil, err
	}
	// insert ke user_identities
	var userIdentity domain.UserIdentity
	createUserIdentityQuery := `
		INSERT INTO user_identities (id, user_id, provider, provider_id, email, password_hash)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, provider, provider_id, email, phone, password_hash, verified, created_at, updated_at
	`
	err = tx.QueryRowxContext(ctx,
		createUserIdentityQuery,
		req.IdIdentity,
		req.ID,
		"local",
		req.Email, // provider_id (pakai email untuk local)
		req.Email,
		req.Password,
	).StructScan(&userIdentity)
	if err != nil {
		return nil, err
	}

	// commit trx
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// response
	return &domain.UserWithIdentity{
		User:         user,
		UserIdentity: userIdentity,
	}, nil
}
