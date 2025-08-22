package handler

import (
	"context"
	authpb "eventify/auth/api"
	"eventify/auth/internal/service"
	"eventify/common/logger"
	"go.uber.org/zap"
)

// AuthGRPCServer provides gRPC endpoints backed by AuthService.
type AuthGRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	service *service.AuthService
}

func NewAuthGRPCServer(s *service.AuthService) *AuthGRPCServer {
	return &AuthGRPCServer{service: s}
}

// Register handles user registration via gRPC.
func (s *AuthGRPCServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	modelReq := toRegisterModel(req)

	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := s.service.RegisterUser(ctx, modelReq); err != nil {
		log.Error(ctx, "register user failed", zap.Error(err))
		return nil, err
	}
	return &authpb.RegisterResponse{Message: "User registered"}, nil
}

// Login handles user login and returns a JWT token.
func (s *AuthGRPCServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	modelReq := toLoginModel(req)

	log := logger.GetOrCreateLoggerFromCtx(ctx)

	token, err := s.service.LoginUser(ctx, modelReq)
	if err != nil {
		log.Error(ctx, "login user failed", zap.Error(err))
		return nil, err
	}
	return &authpb.LoginResponse{Token: token}, nil
}
