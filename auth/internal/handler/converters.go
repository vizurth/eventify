package handler

import (
	authpb "eventify/auth/api"
	"eventify/auth/internal/models"
)

func toRegisterModel(req *authpb.RegisterRequest) models.RegisterRequest {
	return models.RegisterRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
		Role:     req.GetRole(),
	}
}

func toLoginModel(req *authpb.LoginRequest) models.LoginRequest {
	// поддерживаем вход по username или email — репозиторий принимает одно поле
	login := req.GetUsername()
	if login == "" {
		login = req.GetEmail()
	}
	return models.LoginRequest{
		Username: login,
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	}
}
