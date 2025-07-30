package main

import (
	"context"
	"eventify/auth/internal/config"
	"eventify/auth/internal/handler"
	"eventify/auth/internal/repository"
	"eventify/auth/internal/service"
	"eventify/common/logger"
	"eventify/common/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	cfg, _ := config.New()

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)

	_ = postgres.WaitForPostgres(ctx, cfg.Postgres, 10, 1*time.Second)

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
