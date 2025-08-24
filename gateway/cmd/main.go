package main

import (
	"context"
	"eventify/common/logger"
	"eventify/gateway/internal/config"
	"eventify/gateway/internal/service"
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

	cfg, _ := config.New()
	fmt.Println(cfg)

	// Создание и запуск gateway сервиса
	gatewayService := service.NewGatewayService(&cfg, log)
	
	go func() {
		log.Info(ctx, "Starting Eventify Gateway Service...")
		if err := gatewayService.Start(); err != nil {
			log.Fatal(ctx, "Failed to start gateway service: ", zap.Error(err))
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info(ctx, "received shutdown signal, stopping gateway")
}
