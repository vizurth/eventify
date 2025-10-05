package app

import (
	"context"
	"eventify/common/grpc/interceptors"
	"eventify/common/logger"
	"eventify/common/postgres"
	uipb "eventify/user-interact/api"
	"eventify/user-interact/internal/config"
	"eventify/user-interact/internal/handler"
	"eventify/user-interact/internal/repository"
	"eventify/user-interact/internal/service"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config *config.Config
	log    *logger.Logger
	pool   *pgxpool.Pool
	server *grpc.Server
}

func New(ctx context.Context, config *config.Config) (*App, error) {
	ctx, _, err := logger.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	log := logger.GetLoggerFromCtx(ctx)

	pool, err := postgres.New(ctx, config.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize postgres: %w", err)
	}

	uiRepo := repository.NewUserInteractionRepository(pool)
	uiService := service.NewUserInteractionService(ctx, uiRepo, config.Kafka)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptors.TimeoutInterceptor((4 * time.Second)),
		),
	)
	grpcHandler := handler.NewUserInteractionHandler(uiService)

	uipb.RegisterUserInteractionServiceServer(grpcServer, grpcHandler)

	return &App{
		config: config,
		log:    log,
		pool:   pool,
		server: grpcServer,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.UserInteract.Port))
	if err != nil {
		a.log.Fatal(ctx, "failed to listen for gRPC", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		a.log.Info(ctx, fmt.Sprintf("gRPC server listening on port %d", a.config.UserInteract.Port))
		if err = a.server.Serve(lis); err != nil {
			log.Fatal(ctx, "gRPC server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	a.log.Info(ctx, "shutting down gRPC server...")
	a.Shutdown(ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	a.server.GracefulStop()
	a.pool.Close()

	a.log.Info(ctx, "gRPC server shutdown successfully")
}
