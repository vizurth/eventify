package main

import (
	"context"
	"eventify/common/logger"
	"eventify/common/postgres"
	uipb "eventify/user-interact/api"
	"eventify/user-interact/internal/config"
	"eventify/user-interact/internal/handler"
	"eventify/user-interact/internal/repository"
	"eventify/user-interact/internal/service"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func main() {
	cfg, _ := config.New()

	ctx := context.Background()
	ctx, _, _ = logger.New(ctx)

	log := logger.GetLoggerFromCtx(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)

	uiRepo := repository.NewUserInteractionRepository(pool)
	uiService := service.NewUserInteractionService(ctx, uiRepo, cfg.Kafka)
	grpcServer := grpc.NewServer()
	grpcHandler := handler.NewUserInteractionHandler(uiService)

	uipb.RegisterUserInteractionServiceServer(grpcServer, grpcHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.UserInteract.Port))
	if err != nil {
		log.Fatal(ctx, "failed to listen for gRPC", zap.Error(err))
	}

	log.Info(ctx, fmt.Sprintf("gRPC server listening on port %d", cfg.UserInteract.Port))
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(ctx, "gRPC server failed", zap.Error(err))
	}
}
