package middleware

import (
	"booking/internal/domain"
	"booking/pkg/security"
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
		} else if c.Get("session") != "" {
			sessionToken := c.Request().Header.Cookie("session")
			session, err := security.GetSession(c.RequestCtx(), m.rdb, string(sessionToken), m.log)
			if err != nil {
				m.log.Error(err, "failed to get session")
				return utils.ErrorResponse(c, domain.ErrUnauthorized, nil)
			}

			c.Locals(domain.SessionCtxKey, session)
		}

		return c.Next()
	}
}
