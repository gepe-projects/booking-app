package domain

import "errors"

// ! JIKA MENAMBAH ERROR BARU, TAMBAHKAN JUGA DI utils/errorResponse.go
// ! buat switch case baru agar tidak keluar raw error nya
// itulah kenapa di layer repository, harus udah balikin error yang sesuai, bukan error raw dari database
var (
	// general
	ErrInternalServerError = errors.New("internal server error")

	// transport layer error
	ErrInvalidRequest = errors.New("invalid request")
	ErrForbiden       = errors.New("forbiden request")
	ErrToomanyrequest = errors.New("too many request")

	// service error
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// jwt error
	ErrInvalidToken = errors.New("invalid token")

	// auth error
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)
