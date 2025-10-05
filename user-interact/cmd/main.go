package main

import (
	"eventify/common/logger"
	"eventify/user-interact/internal/app"
	"eventify/user-interact/internal/config"
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
