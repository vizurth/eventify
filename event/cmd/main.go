package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"vizurth/eventify/common/logger"
	"vizurth/eventify/common/postgres"
	"vizurth/eventify/event/internal/config"
	"vizurth/eventify/event/internal/handler"
	"vizurth/eventify/event/internal/repository"
	"vizurth/eventify/event/internal/service"
)

func main() {
	cfg, _ := config.New()

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)

	err := postgres.Migrate(ctx, cfg.Postgres, cfg.Event.MigrationPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "migration failed", zap.Error(err))
	}

	router := gin.Default()
	eventRepo := repository.NewEventRepository(pool)
	eventService := service.NewEventService(eventRepo)
	eventHandler := handler.NewEventHandler(eventService, router)

	eventHandler.RegisterRoutes()

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Event.Port), router); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service failed to start", zap.Error(err))
		}
	}()

	select {}
}
