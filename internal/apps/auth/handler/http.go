package handler

import (
	"time"

	"booking/internal/domain"
	"booking/internal/server/middleware"
	"booking/pkg/config"
	"booking/pkg/logger"
	"booking/pkg/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/medama-io/go-useragent"
)

type authHandler struct {
	authUsecase domain.AuthUsecase
	mw          *middleware.Middleware
	log         logger.Logger
	config      *config.Config
	userAgent   *useragent.Parser
}

func NewAuthHandler(
	authUsecase domain.AuthUsecase,
	mw *middleware.Middleware,
	log logger.Logger,
	config *config.Config,
) *authHandler {
	ua := useragent.NewParser()
	return &authHandler{
		userAgent:   ua,
		authUsecase: authUsecase,
		mw:          mw,
		log:         log,
		config:      config,
	}
}

func (h *authHandler) RegisterRoutes(r fiber.Router) {
	r.Post("/register", h.register)
	r.Post("/login", h.mw.LoginLimiter(), h.login)
	r.Get("/sessions", h.mw.Auth(), h.getAllActiveSessions)
}

func (h *authHandler) register(c fiber.Ctx) error {
	var req domain.RegisterDTO
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	res, err := h.authUsecase.RegisterUser(c.RequestCtx(), &req)
	if err != nil {
		return utils.ErrorResponse(c, err, err)
	}
	return c.JSON(domain.HttpResponse{
		Success: true,
		Data:    res,
	})
}

func (h *authHandler) login(c fiber.Ctx) error {
	var req domain.LoginDTO

	if parsedDTO := c.Locals("loginDTO"); parsedDTO != nil {
		req = parsedDTO.(domain.LoginDTO)
	} else {
		// fallback kalau middleware gak dipakai
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
	}

	// get user agent
	ua := h.userAgent.Parse(string(c.RequestCtx().UserAgent()))
	req.UserAgent = string(c.RequestCtx().UserAgent())
	req.Device = ua.Device().String()
	req.IpAddress = c.IP()

	res, token, err := h.authUsecase.Login(c.RequestCtx(), &req)
	if err != nil {
		return utils.ErrorResponse(c, err, nil)
	}
	// set cookies
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(h.config.App.AuthSessionTtl),
		HTTPOnly: true,
		Secure:   h.config.App.Env == "prod",
		SameSite: "none",
		Path:     "/",
	})
	return c.JSON(domain.HttpResponse{
		Success: true,
		Data:    res,
	})
}

func (h *authHandler) getAllActiveSessions(c fiber.Ctx) error {
	userId := c.Locals(domain.SessionCtxKey).(*domain.Session).UserID
	sessions, err := h.authUsecase.GetAllActiveSessions(c.RequestCtx(), userId)
	if err != nil {
		return utils.ErrorResponse(c, err, nil)
	}
	return c.JSON(domain.HttpResponse{
		Success: true,
		Data:    sessions,
	})
}
