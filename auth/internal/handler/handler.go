package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vizurth/eventify/auth/internal/models"
	"vizurth/eventify/auth/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
	router  *gin.Engine
}

func NewAuthHandler(service *service.AuthService, router *gin.Engine) *AuthHandler {
	return &AuthHandler{
		service: service,
		router:  router,
	}
}
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.RegisterUser(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.LoginUser(ctx, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) RegisterRoutes() {
	auth := h.router.Group("/auth")
	auth.POST("/register", h.RegisterHandler)
	auth.POST("/login", h.LoginHandler)
}
