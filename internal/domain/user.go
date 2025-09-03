package domain

import (
	"context"
	"time"
)

type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	ImageURL  *string   `db:"image_url"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserWithIdentity struct {
	User         User
	UserIdentity UserIdentity
}

type UserUsecase interface {
	GetByEmail(ctx context.Context, email string) (*UserWithIdentity, error)
	RegisterUser(ctx context.Context, req *RegisterDTO) (res *UserWithIdentity, err error)
	Login(ctx context.Context, req *LoginDTO) (res *UserWithIdentity, token string, err error)
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*UserWithIdentity, error)
	RegisterUser(ctx context.Context, req *RegisterDTO) (*UserWithIdentity, error)
}
