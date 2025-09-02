package middleware

import (
	"fmt"

	"booking/internal/domain"
	"booking/pkg/security"
	"booking/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) LoginLimiter() fiber.Handler {
	return func(c fiber.Ctx) error {
		var dto domain.LoginDTO

		if err := c.Bind().Body(&dto); err != nil {
			m.log.Error(err, "failed to bind login dto")
			return err // Error akan di-handle oleh error handler fiber
		}

		// Check ban
		delay, err := security.CheckBan(c.RequestCtx(), m.rdb, dto.Email, m.log)
		if err != nil {
			m.log.Error(err, "failed to check ban in redis")
			return utils.ErrorResponse(c, domain.ErrInternalServerError, nil)
		}

		if delay > 0 {
			data := fmt.Sprintf("too many attempts, please try again after %s", delay.String())
			return utils.ErrorResponse(c, domain.ErrToomanyrequest, data)
		}

		return c.Next()
	}
}
