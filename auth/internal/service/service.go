package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
	"vizurth/eventify/auth/internal/models"
	"vizurth/eventify/auth/internal/repository"
	"vizurth/eventify/common/jwt"
)

type AuthService struct {
	repo   *repository.AuthRepository
	secret string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func NewAuthService(repo *repository.AuthRepository, secret string) *AuthService {
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
		return err
	}

	return s.repo.CreateUser(ctx, req.Username, req.Email, hash, req.Role)
}

func (s *AuthService) LoginUser(ctx context.Context, req models.LoginRequest) (string, error) {
	var hashedPassword, role string
	var userId int

	if err := s.repo.GetUser(ctx, req.Username, &hashedPassword, &userId, &role); err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(s.secret, userId, req.Username, req.Email, role, time.Hour*24)
	if err != nil {
		return "", err
	}
	return token, nil

}
