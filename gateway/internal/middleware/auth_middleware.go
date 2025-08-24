package middleware

import (
	"context"
	"eventify/common/jwt"
	"eventify/common/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type AuthMiddleware struct {
	secretKey string
	log       *logger.Logger
}

func NewAuthMiddleware(secretKey string, log *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
		log:       log,
	}
}

func (m *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Проверяем, что это не auth endpoint
		if strings.HasPrefix(r.URL.Path, "/auth/") {
			next.ServeHTTP(w, r)
			return
		}

		// Извлекаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Warn(ctx, "missing authorization header", zap.String("path", r.URL.Path))
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Проверяем формат заголовка (Bearer <token>)
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			m.log.Warn(ctx, "invalid authorization header format", zap.String("header", authHeader))
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]

		// Валидируем JWT токен
		claims, err := m.validateToken(token)
		if err != nil {
			m.log.Warn(ctx, "invalid token", zap.Error(err), zap.String("path", r.URL.Path))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем claims в контекст для использования в handlers
		ctxWithClaims := context.WithValue(ctx, "user_claims", claims)
		rWithClaims := r.WithContext(ctxWithClaims)

		m.log.Debug(ctx, "token validated successfully",
			zap.String("user_id", strconv.Itoa(claims.UserID)),
			zap.String("username", claims.Username),
			zap.String("path", r.URL.Path))

		next.ServeHTTP(w, rWithClaims)
	})
}

type UserClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (m *AuthMiddleware) validateToken(tokenString string) (*UserClaims, error) {
	if tokenString == "" {
		return nil, jwt.ErrInvalidToken
	}

	// Используем существующую JWT реализацию
	claims, err := jwt.ParseToken(tokenString, []byte(m.secretKey))
	if err != nil {
		return nil, err
	}

	return &UserClaims{
		UserID:   claims.UserId,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}, nil
}
