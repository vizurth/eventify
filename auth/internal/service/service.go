package service

import (
	"context"
	"errors"
	"eventify/auth/internal/models"
	"eventify/auth/internal/repository"
	"eventify/common/jwt"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrRefreshExpired    = errors.New("refresh token expired")
)

const (
	accessTokenTTL  = 24 * time.Hour
	refreshTokenTTL = 7 * 24 * time.Hour
	newAccessTTL    = 15 * time.Minute
)

type AuthService struct {
	repo   repository.Repository
	secret string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), nil
}

func NewAuthService(repo repository.Repository, secret string) Service {
	return &AuthService{
		repo:   repo,
		secret: secret,
	}
}

// RegisterUser service register new user
func (s *AuthService) RegisterUser(ctx context.Context, req models.RegisterRequest) error {
	exists, err := s.repo.UserExists(ctx, req.Username, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserAlreadyExists
	}
	hash, err := HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	return s.repo.CreateUser(ctx, req.Username, req.Email, hash, req.Role)
}

// LoginUser service login user in system and get tokens from repo
func (s *AuthService) LoginUser(ctx context.Context, req models.LoginRequest) (string, string, error) {
	var hashedPassword, role string
	var userId int

	if err := s.repo.GetUser(ctx, req.Username, &hashedPassword, &userId, &role); err != nil {
		return "", "", fmt.Errorf("login user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return "", "", fmt.Errorf("login user: invalid credentials")
	}

	accessToken, err := jwt.GenerateToken(s.secret, userId, req.Username, req.Email, role, accessTokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("login user: generate token: %w", err)
	}

	refreshToken := uuid.New().String()
	expiresAt := time.Now().Add(refreshTokenTTL)

	if err := s.repo.SaveRefreshToken(ctx, userId, refreshToken, expiresAt); err != nil {
		return "", "", fmt.Errorf("login user: save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// Logout service logout from system delete refresh-token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.DeleteRefreshToken(ctx, refreshToken)
}

// RefreshToken service refresh access token if time is expired
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	userId, expiresAt, err := s.repo.GetRefreshTokenInfo(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh token info: %w", err)
	}
	if time.Now().After(expiresAt) {
		_ = s.repo.DeleteRefreshToken(ctx, refreshToken)
		return "", ErrRefreshExpired
	}

	var username, email, role string
	if err = s.repo.GetUserForGenerateNewToken(ctx, userId, &username, &email, &role); err != nil {
		return "", fmt.Errorf("refresh token info: %w", err)
	}

	return jwt.GenerateToken(s.secret, userId, username, email, role, newAccessTTL)
}
