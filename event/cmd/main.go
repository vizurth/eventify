package main

import (
	"context"
	"eventify/common/kafka"
	"eventify/common/logger"
	"eventify/common/postgres"
	"eventify/event/internal/config"
	"eventify/event/internal/handler"
	"eventify/event/internal/repository"
	"eventify/event/internal/service"
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

	err := postgres.Migrate(ctx, cfg.Postgres, cfg.Event.MigrationPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "migration failed", zap.Error(err))
	}

	router := gin.Default()
	eventRepo := repository.NewEventRepository(pool)

	producer := kafka.NewProducer([]string{"kafka:9092"}, "events")
	defer producer.Close()

	eventService := service.NewEventService(eventRepo, producer)
	eventHandler := handler.NewEventHandler(eventService, router)

	eventHandler.RegisterRoutes()

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Event.Port), router); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service failed to start", zap.Error(err))
		}
	}()

	select {}
}
