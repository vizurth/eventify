package main

import (
	"context"
	"eventify/common/logger"
	"eventify/notification/internal/app"
	"eventify/notification/internal/config"
	"go.uber.org/zap"
	"os"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx).Error(ctx, "configuration initialization failed", zap.Error(err))
	}

	application, err := app.New(ctx, cfg)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx).Error(ctx, "application initialization failed", zap.Error(err))
		os.Exit(1)
	}

	application.Run(ctx)
}
