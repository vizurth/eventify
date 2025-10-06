package service

import (
	"context"
	"eventify/auth/internal/models"
)

type Service interface {
	RegisterUser(ctx context.Context, req models.RegisterRequest) error
	LoginUser(ctx context.Context, req models.LoginRequest) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}
