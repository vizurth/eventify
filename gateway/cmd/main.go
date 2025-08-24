package main

import (
	"context"
	"eventify/common/logger"
	"eventify/gateway/internal/config"
	"eventify/gateway/internal/handler"
	"eventify/gateway/internal/middleware"
	"fmt"
	"go.uber.org/zap"
	"net/http"
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

	// Создаем gateway handler
	gatewayHandler := handler.NewGatewayHandler(&cfg, log)

	// Создаем middleware для аутентификации
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.SecretKey, log)

	// Настраиваем роуты
	mux := http.NewServeMux()

	// Публичные роуты (без аутентификации)
	mux.HandleFunc("/auth/", gatewayHandler.HandleAuth)

	// Защищенные роуты (с аутентификацией)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/events/", gatewayHandler.HandleEvents)
	protectedMux.HandleFunc("/user-interact/", gatewayHandler.HandleUserInteraction)
	protectedMux.HandleFunc("/registration/", gatewayHandler.HandleRegistration)

	// Применяем middleware аутентификации к защищенным роутам
	mux.Handle("/", authMiddleware.AuthMiddleware(protectedMux))

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Gateway.Port),
		Handler: mux,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Info(ctx, "starting gateway server", zap.Int("port", cfg.Gateway.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(ctx, "gateway server error", zap.Error(err))
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info(ctx, "received shutdown signal, stopping gateway")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Gateway.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error(ctx, "error during server shutdown", zap.Error(err))
	}

	log.Info(ctx, "gateway server stopped")
}
