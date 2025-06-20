package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/vizurth/eventify/internal/authservice"
	"github.com/vizurth/eventify/internal/config"
	"github.com/vizurth/eventify/internal/eventservice"
	"github.com/vizurth/eventify/pkg/logger"
	"github.com/vizurth/eventify/pkg/postgres"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	// создаем конфиг
	cfg, _ := config.NewConfig()

	// прокидываем контект для logger
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	// прокидываем базу данных
	pool, _ := postgres.New(ctx, cfg.PostgresConfig)

	// добавляем authRouter для authService
	authRouter := gin.Default()
	authServ := authservice.NewAuthService(pool, authRouter, []byte(cfg.SecretKey))
	authServ.RegisterRoutes()

	// добавляем eventRouter для eventService
	eventRouter := gin.Default()
	eventServ := eventservice.NewEventService(pool, eventRouter)
	eventServ.RegisterRoutes()

	// запускаем сервисы на разных портах при каких то ошибках будем видет подробную информацию
	//logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service:", zap.Error(http.ListenAndServe(":8081", authRouter)))
	//logger.GetLoggerFromCtx(ctx).Fatal(ctx, "event service:", zap.Error(http.ListenAndServe(":8082", eventRouter)))

	// так как у нас работают сервисы параллельно нужно запускать их по отдельности через горутину
	go func() {
		if err := http.ListenAndServe(":8081", authRouter); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service failed", zap.Error(err))
		}
	}()

	go func() {
		if err := http.ListenAndServe(":8082", eventRouter); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "event service failed", zap.Error(err))
		}
	}()

	// Блокируем main, чтобы горутины не завершились
	select {}

}
