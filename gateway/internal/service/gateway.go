package service

import (
	"context"
	"eventify/common/logger"
	"eventify/gateway/internal/middleware"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"strconv"
	"strings"

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

	// Создаем gRPC соединения к сервисам
	authConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", g.config.Auth.Host, g.config.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	eventConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", g.config.Event.Host, g.config.Event.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to event service: %v", err)
	}
	defer eventConn.Close()

	userInteractConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", g.config.UserInteract.Host, g.config.UserInteract.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to user-interact service: %v", err)
	}
	defer userInteractConn.Close()

	// Создаем gRPC-Gateway мультиплексор
	gwmux := runtime.NewServeMux()

	// Регистрируем сервисы
	if err := authpb.RegisterAuthServiceHandler(ctx, gwmux, authConn); err != nil {
		return fmt.Errorf("failed to register auth service: %v", err)
	}

	if err := eventpb.RegisterEventServiceHandler(ctx, gwmux, eventConn); err != nil {
		return fmt.Errorf("failed to register event service: %v", err)
	}

	if err := uipb.RegisterUserInteractionServiceHandler(ctx, gwmux, userInteractConn); err != nil {
		return fmt.Errorf("failed to register user-interact service: %v", err)
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
	mux.Handle("/", c.Handler(middleware.AuthMiddleware(gwmux, g.logger)))

	// Запускаем сервер
	addr := fmt.Sprintf(":%d", g.config.Server.Port)
	g.logger.Info(ctx, "Starting gateway server on port %d", zap.String("Port: ", strconv.Itoa(g.config.Server.Port)))

	return http.ListenAndServe(addr, mux)
}

func (g *GatewayService) customErrorHandler(ctx context.Context, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	g.logger.Error(ctx, "Gateway error: %v", zap.Error(err))

	// Определяем HTTP статус код на основе ошибки
	status := http.StatusInternalServerError
	if strings.Contains(err.Error(), "not found") {
		status = http.StatusNotFound
	} else if strings.Contains(err.Error(), "invalid") {
		status = http.StatusBadRequest
	} else if strings.Contains(err.Error(), "unauthorized") {
		status = http.StatusUnauthorized
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]interface{}{
		"error": err.Error(),
		"code":  status,
	}

	response, _ := marshaler.Marshal(errorResponse)
	w.Write(response)
}
