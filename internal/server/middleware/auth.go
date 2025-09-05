package middleware

import (
	"time"

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
			session, refreshed, err := m.security.GetSession(c.RequestCtx(), string(sessionToken))
			if err != nil {
				return utils.ErrorResponse(c, domain.ErrUnauthorized, nil)
			}
			if refreshed {
				c.Cookie(&fiber.Cookie{
					Name:     "session",
					Value:    string(sessionToken),
					Expires:  time.Now().Add(m.config.App.AuthSessionTtl),
					HTTPOnly: true,
					Secure:   m.config.App.Env == "prod",
					SameSite: "Lax",
					Path:     "/",
				})
			}
			c.Locals(domain.SessionCtxKey, session)
		}

		return c.Next()
	}
}
