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
	"eventify/common/redis"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	cfg, _ := config.New()

	ctx := context.Background()
	ctx, _, _ = logger.New(ctx)

	log := logger.GetLoggerFromCtx(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)
	redisClient, _ := redis.NewClient(ctx, cfg.Redis)

	authRepo := repository.NewAuthRepository(pool, redisClient)
	authService := service.NewAuthService(authRepo, cfg.Auth.SecretKey)
	grpcServer := grpc.NewServer()
	grpcHandler := handler.NewAuthGRPCServer(authService)

	// Register gRPC server
	authpb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	// Start gRPC server on a dedicated port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Auth.Port))
	if err != nil {
		log.Fatal(ctx, "failed to listen for gRPC", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info(ctx, fmt.Sprintf("gRPC server listening on port %d", cfg.Auth.Port))
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatal(ctx, "gRPC server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info(ctx, "shutting down gRPC server...")

	grpcServer.GracefulStop()
	pool.Close()

	if err := redisClient.Close(); err != nil {
		log.Fatal(ctx, "failed to close redis client", zap.Error(err))
	}

	log.Info(ctx, "successfully shutdown gRPC server")
}
