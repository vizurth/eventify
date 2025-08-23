package main

import (
	"context"
	authpb "eventify/auth/api"
	"eventify/auth/internal/config"
	"eventify/auth/internal/handler"
	"eventify/auth/internal/repository"
	"eventify/auth/internal/service"
	"eventify/common/logger"
	"eventify/common/postgres"
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

	if err := postgres.Migrate(ctx, cfg.Postgres, cfg.Auth.MigrationPath); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "migration failed", zap.Error(err))
	}

	authRepo := repository.NewAuthRepository(pool)
	authService := service.NewAuthService(authRepo, cfg.Auth.SecretKey)
	grpcServer := grpc.NewServer()
	grpcHandler := handler.NewAuthGRPCServer(authService)

	// Register gRPC server
	authpb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	// Start gRPC server on a dedicated port
	grpcPort := 9091
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to listen for gRPC", zap.Error(err))
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "gRPC server failed", zap.Error(err))
		}
	}()

	// Setup and start grpc-gateway on HTTP port from config
	mux := gw.NewServeMux()
	connOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), connOpts); err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to register grpc-gateway", zap.Error(err))
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Auth.Port),
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "http gateway failed", zap.Error(err))
		}
	}()

	select {}
} 