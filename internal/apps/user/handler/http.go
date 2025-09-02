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
	r.Get("/", h.getUser)
}

func (h *userHandler) getUser(c fiber.Ctx) error {
	res, err := h.UseCase.GetByEmail(c.RequestCtx(), c.Query("email"))
	if err != nil {
		return utils.ErrorResponse(c, err, nil)
	}

	return c.JSON(fiber.Map{
		"data": res,
	})
}
