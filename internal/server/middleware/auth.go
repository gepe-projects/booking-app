package middleware

import (
	"booking/internal/domain"
	"booking/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) Auth() fiber.Handler {
	return func(c fiber.Ctx) error {
		// var token string

		// bisa pakai jwt atau session
		if c.Get("Authorization") != "" {
			m.log.Error(nil, "JWT NOT IMPLEMENTED YET")
			return c.Status(fiber.StatusInternalServerError).JSON(domain.HttpResponse{
				Success: false,
				Message: domain.ErrInternalServerError,
			})
		} else if string(c.Request().Header.Cookie("session")) != "" {
			sessionToken := c.Request().Header.Cookie("session")
			session, err := m.security.GetSession(c.RequestCtx(), string(sessionToken))
			if err != nil {
				return utils.ErrorResponse(c, domain.ErrUnauthorized, nil)
			}
			c.Locals(domain.SessionCtxKey, session)
		}

		return c.Next()
	}
}
