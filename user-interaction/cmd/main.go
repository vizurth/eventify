package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"vizurth/eventify/common/logger"
	"vizurth/eventify/common/postgres"
	"vizurth/eventify/user-interaction/internal/config"
	"vizurth/eventify/user-interaction/internal/handler"
	"vizurth/eventify/user-interaction/internal/repository"
	"vizurth/eventify/user-interaction/internal/service"
)

func main() {
	cfg, _ := config.New()

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)

	err := postgres.Migrate(ctx, cfg.Postgres, cfg.UserInteract.MigrationPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "migration failed", zap.Error(err))
	}

	router := gin.Default()
	eventRepo := repository.NewUserInteractionRepository(pool)
	eventService := service.NewUserInteractionService(eventRepo)
	eventHandler := handler.NewUserInteractionHandler(eventService, router, []byte(cfg.UserInteract.Secret))

	eventHandler.RegisterRoutes()

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.UserInteract.Port), router); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service failed to start", zap.Error(err))
		}
	}()

	select {}
}
