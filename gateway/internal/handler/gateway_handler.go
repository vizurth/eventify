package handler

import (
	"context"
	"eventify/common/logger"
	"eventify/gateway/internal/config"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type GatewayHandler struct {
	config *config.Config
	log    *logger.Logger
	client *http.Client
}

func NewGatewayHandler(cfg *config.Config, log *logger.Logger) *GatewayHandler {
	return &GatewayHandler{
		config: cfg,
		log:    log,
		client: &http.Client{},
	}
}

// HandleAuth обрабатывает запросы к auth сервису
func (h *GatewayHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	// Убираем префикс /auth из пути
	path := strings.TrimPrefix(r.URL.Path, "/auth")
	if path == "" {
		path = "/"
	}

	// Определяем endpoint
	var endpoint string
	switch {
	case r.Method == "POST" && path == "/register":
		endpoint = "/auth/register"
	case r.Method == "POST" && path == "/login":
		endpoint = "/auth/login"
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.forwardRequest(w, r, h.config.Auth.URL+endpoint)
}

// HandleEvents обрабатывает запросы к event сервису
func (h *GatewayHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	// Убираем префикс /events из пути
	path := strings.TrimPrefix(r.URL.Path, "/events")
	if path == "" {
		path = "/"
	}

	// Определяем endpoint
	var endpoint string
	switch {
	case r.Method == "GET" && path == "/":
		endpoint = "/events/"
	case r.Method == "POST" && path == "/":
		endpoint = "/events/"
	case r.Method == "GET" && strings.HasPrefix(path, "/"):
		// Извлекаем ID события
		params := strings.TrimPrefix(path, "/")
		endpoint = fmt.Sprintf("/events/%s", params)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.forwardRequest(w, r, h.config.Event.URL+endpoint)
}

// HandleUserInteraction обрабатывает запросы к user-interaction сервису
func (h *GatewayHandler) HandleUserInteraction(w http.ResponseWriter, r *http.Request) {
	// Убираем префикс /user-interact из пути
	path := strings.TrimPrefix(r.URL.Path, "/user-interact")
	if path == "" {
		path = "/"
	}

	// Определяем endpoint
	var endpoint string
	switch {
	case r.Method == "POST" && path == "/":
		endpoint = "/user-interact/"
	case r.Method == "GET" && strings.HasPrefix(path, "/event/"):
		eventID := strings.TrimPrefix(path, "/event/")
		endpoint = fmt.Sprintf("/user-interact/event/%s", eventID)
	case r.Method == "PUT" && strings.HasPrefix(path, "/"):
		reviewID := strings.TrimPrefix(path, "/")
		endpoint = fmt.Sprintf("/user-interact/%s", reviewID)
	case r.Method == "DELETE" && strings.HasPrefix(path, "/"):
		reviewID := strings.TrimPrefix(path, "/")
		endpoint = fmt.Sprintf("/user-interact/%s", reviewID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.forwardRequest(w, r, h.config.UserInteraction.URL+endpoint)
}

// HandleRegistration обрабатывает запросы к registration endpoints
func (h *GatewayHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	// Убираем префикс /registration из пути
	path := strings.TrimPrefix(r.URL.Path, "/registration")
	if path == "" {
		path = "/"
	}

	// Определяем endpoint
	var endpoint string
	switch {
	case r.Method == "POST" && strings.HasPrefix(path, "/event/"):
		eventID := strings.TrimPrefix(path, "/event/")
		endpoint = fmt.Sprintf("/registration/event/%s", eventID)
	case r.Method == "DELETE" && strings.HasPrefix(path, "/event/"):
		eventID := strings.TrimPrefix(path, "/event/")
		endpoint = fmt.Sprintf("/registration/event/%s", eventID)
	case r.Method == "GET" && strings.HasPrefix(path, "/event/"):
		eventID := strings.TrimPrefix(path, "/event/")
		endpoint = fmt.Sprintf("/registration/event/%s", eventID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.forwardRequest(w, r, h.config.UserInteraction.URL+endpoint)
}

// forwardRequest пересылает HTTP запрос к указанному сервису
func (h *GatewayHandler) forwardRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	ctx := r.Context()

	h.log.Debug(ctx, "forwarding request",
		zap.String("method", r.Method),
		zap.String("original_path", r.URL.Path),
		zap.String("target_url", targetURL))

	// Создаем новый запрос
	req, err := http.NewRequestWithContext(ctx, r.Method, targetURL, r.Body)
	if err != nil {
		h.log.Error(ctx, "failed to create request", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Добавляем заголовки для логирования
	req.Header.Set("X-Forwarded-For", r.RemoteAddr)
	req.Header.Set("X-Forwarded-Proto", r.URL.Scheme)
	req.Header.Set("X-Forwarded-Host", r.Host)

	// Выполняем запрос
	resp, err := h.client.Do(req)
	if err != nil {
		h.log.Error(ctx, "failed to forward request", zap.Error(err), zap.String("target", targetURL))
		http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Копируем заголовки ответа
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Устанавливаем статус код
	w.WriteHeader(resp.StatusCode)

	// Копируем тело ответа
	if _, err := io.Copy(w, resp.Body); err != nil {
		h.log.Error(ctx, "failed to copy response body", zap.Error(err))
	}

	h.log.Debug(ctx, "request forwarded successfully",
		zap.Int("status", resp.StatusCode),
		zap.String("path", r.URL.Path))
}

// GetUserClaims извлекает claims пользователя из контекста
func (h *GatewayHandler) GetUserClaims(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value("user_claims").(*UserClaims)
	return claims, ok
}

type UserClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
