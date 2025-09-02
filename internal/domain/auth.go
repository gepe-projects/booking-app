package domain

import "context"

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email" message:"Valid email is required"`
	Password string `json:"password" validate:"required,min=6,max=150" message:"Password is required and minimum length is 6"`

	// device info
	Device    string `json:"device"`
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

type RegisterDTO struct {
	ID       string
	Name     string `json:"name" validate:"required,min=2,max=100" message:"Name is required and minimum length is 2"`
	ImageURL string
	Role     string
	// identity
	IdIdentity string
	Email      string `json:"email" validate:"required,email" message:"Valid email is required"`
	Password   string `json:"password" validate:"required,min=6,max=150" message:"Password is required and minimum length is 6"`
}

type AuthUsecase interface {
	RegisterUser(ctx context.Context, req *RegisterDTO) (res *UserWithIdentity, err error)
	Login(ctx context.Context, req *LoginDTO) (res *UserWithIdentity, token string, err error)
}
