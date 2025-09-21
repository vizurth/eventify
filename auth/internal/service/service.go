package service

import (
	"context"
	"errors"
	"eventify/auth/internal/models"
	"eventify/auth/internal/repository"
	"eventify/common/jwt"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var HashError = errors.New("hash error")

type AuthService struct {
	repo   repository.Repository
	secret string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), fmt.Errorf("hash error: %w", err)
}

func NewAuthService(repo repository.Repository, secret string) *AuthService {
	return &AuthService{
		repo:   repo,
		secret: secret,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, req models.RegisterRequest) error {
	exitst, err := s.repo.UserExists(ctx, req.Username, req.Email)
	if err != nil {
		return err
	}
	if exitst {
		return errors.New("user already exists")
	}
	hash, err := HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	return s.repo.CreateUser(ctx, req.Username, req.Email, hash, req.Role)
}

func (s *AuthService) LoginUser(ctx context.Context, req models.LoginRequest) (string, error) {
	var hashedPassword, role string
	var userId int

	if err := s.repo.GetUser(ctx, req.Username, &hashedPassword, &userId, &role); err != nil {
		return "", fmt.Errorf("login user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("login user: %w", err)
	}

	token, err := jwt.GenerateToken(s.secret, userId, req.Username, req.Email, role, time.Hour*24)
	if err != nil {
		return "", fmt.Errorf("login user: generate token: %w", err)
	}
	return token, nil

}
