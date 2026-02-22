package service

import (
	"context"
	"eventify/common/logger"
	"eventify/gateway/internal/middleware"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "eventify/auth/api"
	eventpb "eventify/event/api"
	"eventify/gateway/internal/config"
	uipb "eventify/user-interact/api"
)

type GatewayService struct {
	config *config.Config
	logger *logger.Logger
}

func NewGatewayService(cfg *config.Config, logger *logger.Logger) *GatewayService {
	return &GatewayService{
		config: cfg,
		logger: logger,
	}
}

func (g *GatewayService) Start() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Создаем gRPC соединения к сервисам
	authConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", g.config.Auth.Host, g.config.Auth.Port),
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to auth service: %w", err)
	}
	defer func() { _ = authConn.Close() }()

	eventConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", g.config.Event.Host, g.config.Event.Port),
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to event service: %w", err)
	}
	defer func() { _ = eventConn.Close() }()

	userInteractConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", g.config.UserInteract.Host, g.config.UserInteract.Port),
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to user-interact service: %w", err)
	}
	defer func() { _ = userInteractConn.Close() }()

	// Создаем gRPC-Gateway мультиплексор
	gwmux := runtime.NewServeMux()

	// Регистрируем сервисы
	if err := authpb.RegisterAuthServiceHandler(ctx, gwmux, authConn); err != nil {
		return fmt.Errorf("failed to register auth service: %w", err)
	}

	if err := eventpb.RegisterEventServiceHandler(ctx, gwmux, eventConn); err != nil {
		return fmt.Errorf("failed to register event service: %w", err)
	}

	if err := uipb.RegisterUserInteractionServiceHandler(ctx, gwmux, userInteractConn); err != nil {
		return fmt.Errorf("failed to register user-interact service: %w", err)
	}

	// Создаем HTTP сервер
	mux := http.NewServeMux()

	// Оборачиваем gRPC-Gateway через CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:9098"}, // порт Swagger UI
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Регистрируем маршруты с CORS и JWT middleware
	handler := middleware.AuthMiddleware(gwmux, g.logger)     // JWT middleware
	handler = middleware.LoggingMiddleware(handler, g.logger) // Логирование
	handler = c.Handler(handler)                              // CORS

	mux.Handle("/", handler)

	// Запускаем сервер
	addr := fmt.Sprintf(":%d", g.config.Server.Port)
	g.logger.Info(ctx, "starting gateway server", zap.Int("port", g.config.Server.Port))

	return http.ListenAndServe(addr, mux)
}

func (g *GatewayService) customErrorHandler(ctx context.Context, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	g.logger.Error(ctx, "gateway error", zap.Error(err))

	// Определяем HTTP статус код на основе ошибки
	httpStatus := http.StatusInternalServerError
	if strings.Contains(err.Error(), "not found") {
		httpStatus = http.StatusNotFound
	} else if strings.Contains(err.Error(), "invalid") {
		httpStatus = http.StatusBadRequest
	} else if strings.Contains(err.Error(), "unauthorized") {
		httpStatus = http.StatusUnauthorized
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	errorResponse := map[string]interface{}{
		"error": err.Error(),
		"code":  httpStatus,
	}

	response, _ := marshaler.Marshal(errorResponse)
	_, _ = w.Write(response)
}
