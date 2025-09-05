package handler

import (
	"booking/internal/domain"
	"booking/internal/server/middleware"
	"booking/pkg/logger"
	"booking/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

type userHandler struct {
	UseCase    domain.UserUsecase
	middleware *middleware.Middleware
	log        logger.Logger
}

func NewUserHandler(useCase domain.UserUsecase, middleware *middleware.Middleware, logger logger.Logger) *userHandler {
	return &userHandler{UseCase: useCase, middleware: middleware, log: logger}
}

func (h *userHandler) RegisterRoutes(r fiber.Router) {
	r.Use(h.middleware.Auth())
	r.Get("/me", h.getUser)
}

func (h *userHandler) getUser(c fiber.Ctx) error {
	val := c.Locals(domain.SessionCtxKey)
	session, ok := val.(*domain.Session)
	if !ok || session == nil {
		return utils.ErrorResponse(c, domain.ErrUnauthorized, nil)
	}

	c.Response().Header.Set("Cache-Control", "private, max-age=60")

	return c.JSON(fiber.Map{
		"data": session,
	})
}
