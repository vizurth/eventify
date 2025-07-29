package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"vizurth/eventify/auth/internal/config"
	"vizurth/eventify/auth/internal/handler"
	"vizurth/eventify/auth/internal/repository"
	"vizurth/eventify/auth/internal/service"
	"vizurth/eventify/common/logger"
	"vizurth/eventify/common/postgres"
)

func main() {
	cfg, _ := config.New()

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)

	err := postgres.Migrate(ctx, cfg.Postgres, cfg.Auth.MigrationPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "migration failed", zap.Error(err))
	}

	router := gin.Default()
	authRepo := repository.NewAuthRepository(pool)
	authService := service.NewAuthService(authRepo, cfg.Auth.SecretKey)
	authHandler := handler.NewAuthHandler(authService, router)

	authHandler.RegisterRoutes()

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Auth.Port), router); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service failed to start", zap.Error(err))
		}
	}()

	select {}
}
