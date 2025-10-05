package main

import (
	"context"
	"eventify/common/logger"
	"eventify/event/internal/app"
	"eventify/event/internal/config"
	"go.uber.org/zap"
	"os"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx).Fatal(ctx, "failed to load configuration", zap.Error(err))
	}

	application, err := app.New(ctx, cfg)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx).Fatal(ctx, "failed to initialize application", zap.Error(err))
		os.Exit(1)
	}

	application.Run(ctx)

}
