package main

import (
	"context"
	uipb "eventify/user-interaction/api"
	"eventify/common/logger"
	"eventify/common/postgres"
	"eventify/user-interaction/internal/config"
	"eventify/user-interaction/internal/handler"
	"eventify/user-interaction/internal/repository"
	"eventify/user-interaction/internal/service"
	"fmt"
	"net"
	"net/http"
	"time"

	gw "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, _ := config.New()
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)
	_ = postgres.WaitForPostgres(ctx, cfg.Postgres, 10, 1*time.Second)

	if err := postgres.Migrate(ctx, cfg.Postgres, cfg.UserInteract.MigrationPath); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "migration failed", zap.Error(err))
	}

	repo := repository.NewUserInteractionRepository(pool)
	service := service.NewUserInteractionService(repo)
	rpcHandler := handler.NewUserInteractionHandler(service)

	grpcServer := grpc.NewServer()
	uipb.RegisterUserInteractionServiceServer(grpcServer, rpcHandler)

	grpcPort := 9093
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to listen for gRPC", zap.Error(err))
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "gRPC server failed", zap.Error(err))
		}
	}()

	mux := gw.NewServeMux()
	connOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := uipb.RegisterUserInteractionServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), connOpts); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to register grpc-gateway", zap.Error(err))
	}

	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", cfg.UserInteract.Port), Handler: mux}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "http gateway failed", zap.Error(err))
		}
	}()

	select {}
}
