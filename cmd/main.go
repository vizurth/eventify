package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/vizurth/eventify/internal/authservice"
	"github.com/vizurth/eventify/internal/config"
	"github.com/vizurth/eventify/pkg/logger"
	"github.com/vizurth/eventify/pkg/postgres"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg, _ := config.NewConfig()

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	pool, _ := postgres.New(ctx, cfg.PostgresConfig)

	authRouter := gin.Default()
	authServ := authservice.NewAuthService(pool, authRouter, []byte(cfg.SecretKey))
	authServ.RegisterRoutes()

	logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service:", zap.Error(http.ListenAndServe(":8081", authRouter)))
}
