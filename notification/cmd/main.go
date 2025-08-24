package main

import (
	"context"
	"eventify/common/logger"
	"eventify/notification/internal/config"
	"eventify/notification/internal/service"
	"eventify/notification/internal/wsserver"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	ctx, _, _ = logger.New(ctx)
	log := logger.GetLoggerFromCtx(ctx)
	// Загружаем конфигурацию
	cfg, err := config.New()
	if err != nil {
		log.Fatal(ctx, "failed to load config", zap.Error(err))
	}

	// Настраиваем Kafka конфигурацию

	// Создаем WebSocket сервер
	wsServer := wsserver.NewWsServer(fmt.Sprintf(":%d", cfg.Notification.Port), log)

	// Создаем notification сервис
	notificationService := service.NewNotificationService(ctx, cfg.Kafka, wsServer, log)

	// Запускаем notification сервис
	if err := notificationService.Start(ctx); err != nil {
		log.Fatal(ctx, "failed to start notification service", zap.Error(err))
	}

	// Запускаем WebSocket сервер в горутине
	go func() {
		log.Info(ctx, "starting websocket server", zap.Int("port", cfg.Notification.Port))
		if err := wsServer.Start(); err != nil {
			log.Error(ctx, "websocket server error", zap.Error(err))
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info(ctx, "received shutdown signal, stopping services")

	// Graceful shutdown
	if err := notificationService.Stop(); err != nil {
		log.Error(ctx, "error stopping notification service", zap.Error(err))
	}

	log.Info(ctx, "notification service stopped")
}
