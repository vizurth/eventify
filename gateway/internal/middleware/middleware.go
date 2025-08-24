package middleware

import (
	"eventify/common/logger"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Разрешаем auth запросы без токена
		if strings.HasPrefix(r.URL.Path, "/auth/") || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем наличие токена для всех остальных запросов
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Warn(ctx, "Unauthorized request to %s: missing Authorization header", zap.String("path:", r.URL.Path))
			return
		}

		// Проверяем формат токена
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Warn(ctx, "Unauthorized request to %s: invalid Authorization header format", zap.String("path:", r.URL.Path))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Warn(ctx, "Unauthorized request to %s: empty token", zap.String("path:", r.URL.Path))
			return
		}

		// Здесь можно добавить дополнительную валидацию токена
		// Например, проверить его через auth сервис
		log.Debug(ctx, "Authorized request to %s with token: %s", zap.String("path:", r.URL.Path), zap.String("token:", token[:10]+"..."))

		// Передаем токен дальше в заголовке для использования в gRPC сервисах
		r.Header.Set("X-User-Token", token)

		next.ServeHTTP(w, r)
	})
}
