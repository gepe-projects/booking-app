package utils

import (
	"errors"

	"booking/internal/domain"

	"github.com/gofiber/fiber/v3"
)

// ! HARUS IMPLEMENTASI SEMUA ERROR YANG ADA DI DOMAIN
func ErrorResponse(c fiber.Ctx, err error, data any) error {
	var statusCode int
	response := domain.HttpResponse{
		Success: false,
		Message: nil,
		Data:    data,
	}

	switch {
	// general
	case errors.Is(err, domain.ErrInternalServerError):
		response.Message = domain.ErrInternalServerError.Error()
		statusCode = fiber.StatusInternalServerError
	// transport layer error
	case errors.Is(err, domain.ErrInvalidRequest):
		response.Message = domain.ErrInvalidRequest.Error()
		statusCode = fiber.StatusBadRequest
	case errors.Is(err, domain.ErrForbiden):
		response.Message = domain.ErrForbiden.Error()
		statusCode = fiber.StatusForbidden
	case errors.Is(err, domain.ErrToomanyrequest):
		response.Message = domain.ErrToomanyrequest.Error()
		statusCode = fiber.StatusTooManyRequests
	// service error
	case errors.Is(err, domain.ErrUserNotFound):
		response.Message = domain.ErrUserNotFound.Error()
		statusCode = fiber.StatusNotFound
	case errors.Is(err, domain.ErrUserAlreadyExists):
		response.Message = domain.ErrUserAlreadyExists.Error()
		statusCode = fiber.StatusConflict
	// jwt error
	case errors.Is(err, domain.ErrInvalidToken):
		response.Message = domain.ErrInvalidToken.Error()
		statusCode = fiber.StatusUnauthorized
	// auth error
	case errors.Is(err, domain.ErrInvalidCredentials):
		response.Message = domain.ErrInvalidCredentials.Error()
		statusCode = fiber.StatusUnauthorized
	case errors.Is(err, domain.ErrUnauthorized):
		response.Message = domain.ErrUnauthorized.Error()
		statusCode = fiber.StatusUnauthorized
	default:
		response.Message = err.Error()
		statusCode = fiber.StatusInternalServerError
	}

	return c.Status(statusCode).JSON(response)
}
