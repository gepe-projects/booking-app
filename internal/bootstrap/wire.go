package bootstrap

import (
	authHandler "booking/internal/apps/auth/handler"
	authUsecase "booking/internal/apps/auth/usecase"
	userHandler "booking/internal/apps/user/handler"
	ur "booking/internal/apps/user/repository"
	userUsecase "booking/internal/apps/user/usecase"
	"booking/internal/server"
	"booking/internal/server/middleware"
	"booking/pkg/config"
	"booking/pkg/database"
	"booking/pkg/logger"
	"booking/pkg/redis"
	"booking/pkg/security"
)

type Apps struct {
	Config *config.Config
	Log    logger.Logger
	Server *server.FiberApp
}

func InitializeApps() *Apps {
	config := ProvideConfig()
	logger := ProvideLogger(config)
	db := database.InitDB(&config.Database, logger)
	rdb := redis.NewClient(&config.Redis, logger)

	// security
	security := security.NewSecurity(rdb, logger)

	// repository
	userRepo := ur.NewUserRepository(db)

	// usecase
	userUsecase := userUsecase.NewUserUseCase(userRepo, security, logger)
	authUsecase := authUsecase.NewAuthUsecase(userUsecase, security, logger)

	// middleware
	middlewares := middleware.NewMiddlewares(security, rdb, logger)

	// handler
	userHandler := userHandler.NewUserHandler(userUsecase, middlewares, logger)
	authHandler := authHandler.NewAuthHandler(authUsecase, middlewares, logger, config)

	// server
	srv := server.NewFiber(&config.Gateway, logger)

	// prefix route
	v1 := srv.App.Group("/api/v1")

	// register routes
	authHandler.RegisterRoutes(v1.Group("/auth"))
	userHandler.RegisterRoutes(v1.Group("/users"))

	return &Apps{
		Config: config,
		Log:    logger,
		Server: srv,
	}
}

func ProvideConfig() *config.Config {
	return config.LoadConfig()
}

func ProvideLogger(cfg *config.Config) logger.Logger {
	return logger.New("booking", cfg)
}
