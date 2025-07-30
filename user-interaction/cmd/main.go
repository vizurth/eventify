package main

import (
	"context"
	"eventify/common/kafka"
	"eventify/common/logger"
	"eventify/common/postgres"
	"eventify/user-interaction/internal/config"
	"eventify/user-interaction/internal/handler"
	"eventify/user-interaction/internal/repository"
	"eventify/user-interaction/internal/service"
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

	err := postgres.Migrate(ctx, cfg.Postgres, cfg.UserInteract.MigrationPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "migration failed", zap.Error(err))
	}

	router := gin.Default()
	eventRepo := repository.NewUserInteractionRepository(pool)

	producer := kafka.NewProducer([]string{"kafka:9092"}, "events")
	defer producer.Close()

	eventService := service.NewUserInteractionService(eventRepo, producer)
	eventHandler := handler.NewUserInteractionHandler(eventService, router, []byte(cfg.UserInteract.Secret))

	eventHandler.RegisterRoutes()

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.UserInteract.Port), router); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service failed to start", zap.Error(err))
		}
	}()

	select {}
}
