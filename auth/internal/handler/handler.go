package handler

import (
	"context"
	authpb "eventify/auth/api"
	"eventify/auth/internal/service"
	"eventify/common/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthGRPCServer provides gRPC endpoints backed by AuthService.
type AuthGRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	service service.Service
}

func NewAuthGRPCServer(s service.Service) *AuthGRPCServer {
	return &AuthGRPCServer{service: s}
}

// Register handles user registration via gRPC.
func (s *AuthGRPCServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	modelReq := toRegisterModel(req)

	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := req.ValidateAll(); err != nil {
		log.Error(ctx, "register handler: validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.service.RegisterUser(ctx, modelReq); err != nil {
		log.Error(ctx, "register handler: register user failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authpb.RegisterResponse{Message: "User registered"}, nil
}

// Login handles user login and returns access and refresh token.
func (s *AuthGRPCServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	modelReq := toLoginModel(req)

	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := req.ValidateAll(); err != nil {
		log.Error(ctx, "login handler: validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	access, refresh, err := s.service.LoginUser(ctx, modelReq)
	if err != nil {
		log.Error(ctx, "login handler: login user failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authpb.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// Logout handles user logout and returns success message.
func (s *AuthGRPCServer) Logout(ctx context.Context, req *authpb.RefreshRequest) (*authpb.LogoutResponse, error) {
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := s.service.Logout(ctx, req.RefreshToken); err != nil {
		log.Error(ctx, "logout handler: failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authpb.LogoutResponse{
		Message: "User logout",
	}, nil
}

// Refresh handles access token refresh.
func (s *AuthGRPCServer) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.LoginResponse, error) {
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := req.ValidateAll(); err != nil {
		log.Error(ctx, "refresh handler: validation failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accessToken, err := s.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		log.Error(ctx, "refresh handler: failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authpb.LoginResponse{
		AccessToken: accessToken,
	}, nil
}
