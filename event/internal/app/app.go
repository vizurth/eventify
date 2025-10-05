package app

import (
	"context"
	"eventify/common/grpc/interceptors"
	"eventify/common/logger"
	"eventify/common/postgres"
	myredis "eventify/common/redis"
	eventpb "eventify/event/api"
	"eventify/event/internal/config"
	"eventify/event/internal/handler"
	"eventify/event/internal/repository"
	"eventify/event/internal/service"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config *config.Config
	log    *logger.Logger
	pool   *pgxpool.Pool
	redis  *redis.Client
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
		return nil, fmt.Errorf("failed to initialize postgres pool: %w", err)
	}

	redisClient, err := myredis.NewClient(ctx, config.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis client: %w", err)
	}

	eventRepo := repository.NewEventRepository(pool, redisClient)
	eventService := service.NewEventService(ctx, eventRepo, config.Kafka)
	grpcHandler := handler.NewEventHandler(eventService)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptors.TimeoutInterceptor((4 * time.Second)),
		),
	)
	eventpb.RegisterEventServiceServer(grpcServer, grpcHandler)

	return &App{
		config: config,
		log:    log,
		pool:   pool,
		redis:  redisClient,
		server: grpcServer,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.Event.Port))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "failed to listen for gRPC", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		a.log.Info(ctx, "gRPC server listening on", zap.Int("port", a.config.Event.Port))
		if err := a.server.Serve(lis); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "gRPC server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	a.log.Info(ctx, "shutting down gRPC server...")

	a.Shutdown(ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	a.server.GracefulStop()
	a.pool.Close()

	if err := a.redis.Close(); err != nil {
		a.log.Fatal(ctx, "failed to close redis client", zap.Error(err))
	}

	a.log.Info(ctx, "gracefully shut down gRPC server...")
}
