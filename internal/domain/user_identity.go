package domain

import "time"

type UserIdentity struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Provider     string    `json:"provider" db:"provider"`
	ProviderID   string    `json:"provider_id" db:"provider_id"`
	Email        *string   `json:"email" db:"email"`
	Phone        *string   `json:"phone,omitempty" db:"phone"`
	PasswordHash *string   `json:"-" db:"password_hash"`
	Verified     bool      `json:"verified" db:"verified"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
