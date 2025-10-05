package app

import (
	"context"
	"eventify/common/logger"
	"eventify/notification/internal/config"
	"eventify/notification/internal/service"
	"eventify/notification/internal/wsserver"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
)

type App struct {
	config  *config.Config
	log     *logger.Logger
	server  wsserver.WsServer
	service *service.NotificationService
}

func New(ctx context.Context, config *config.Config) (*App, error) {
	ctx, _, err := logger.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	log := logger.GetLoggerFromCtx(ctx)

	webSockerServer := wsserver.NewWsServer(fmt.Sprintf(":%d", config.Notification.Port), log)
	notificationService := service.NewNotificationService(ctx, config.Kafka, webSockerServer, log)

	return &App{
		config:  config,
		log:     log,
		server:  webSockerServer,
		service: notificationService,
	}, nil
}

func (a *App) Run(ctx context.Context) {

	if err := a.service.Start(ctx); err != nil {
		log.Fatal(ctx, "failed to start notification service", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		a.log.Info(ctx, "starting websocket server", zap.Int("port", a.config.Notification.Port))
		if err := a.server.Start(); err != nil {
			a.log.Error(ctx, "websocket server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	a.log.Info(ctx, "shutting down websocket server")

	a.Shutdown(ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	if err := a.service.Stop(); err != nil {
		a.log.Error(ctx, "error stopping notification service", zap.Error(err))
	}

	a.log.Info(ctx, "notification service stopped")
}
