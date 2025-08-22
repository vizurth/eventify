package main

import (
	"context"
	"eventify/common/logger"
	"eventify/common/postgres"
	eventpb "eventify/event/api"
	"eventify/event/internal/config"
	"eventify/event/internal/handler"
	"eventify/event/internal/repository"
	"eventify/event/internal/service"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {
	cfg, _ := config.New()

	ctx := context.Background()
	ctx, _, _ = logger.New(ctx)
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)

	eventRepo := repository.NewEventRepository(pool)

	// Kafka producer (topic: events)

	eventService := service.NewEventService(ctx, eventRepo, cfg.Kafka)
	grpcHandler := handler.NewEventHandler(eventService)

	grpcServer := grpc.NewServer()
	eventpb.RegisterEventServiceServer(grpcServer, grpcHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Event.Port))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to listen for gRPC", zap.Error(err))
	}
	log.Info(ctx, "gRPC server listening on", zap.Int("port", cfg.Event.Port))
	if err := grpcServer.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "gRPC server failed", zap.Error(err))
	}
}
