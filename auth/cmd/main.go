package main

import (
	"context"
	"eventify/auth/internal/app"
	"eventify/auth/internal/config"
	"eventify/common/logger"
	"os"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()

	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx).Fatal(ctx, "failed to load config")
	}

	application, err := app.New(ctx, cfg)
	if err != nil {
		logger.GetOrCreateLoggerFromCtx(ctx).Fatal(ctx, "failed to initialize application")
		os.Exit(1)
	}

	application.Run(ctx)
}
