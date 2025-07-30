package handler

import (
	"eventify/common/jwt"
	"eventify/user-interaction/internal/models"
	"eventify/user-interaction/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserInteractionHandler struct {
	service *service.UserInteractionService
	router  *gin.Engine
	secret  []byte
}

func NewUserInteractionHandler(service *service.UserInteractionService, router *gin.Engine, secret []byte) *UserInteractionHandler {
	return &UserInteractionHandler{
		service: service,
		router:  router,
		secret:  secret,
	}
}

// CreateNewReviews создаем новые отзывы
func (h *UserInteractionHandler) CreateNewReviews(c *gin.Context) {
	ctx := c.Request.Context()

	// проверяем корректность запроса
	var req models.ReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// добавляем review в таблицу
	if err := h.service.CreateNewReviews(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review created"})
}

// GetCurrentReviewsByEventID получение review по корректному id
func (h *UserInteractionHandler) GetCurrentReviewsByEventID(c *gin.Context) {
	ctx := c.Request.Context()
	//  получаем eventID
	eventIdParam := c.Param("id")
	eventId, err := strconv.Atoi(eventIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// составляем массив review из таблицы
	var reviews []models.ReviewResp

	if err = h.service.GetCurrentReviewsByEventID(ctx, eventId, &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// UpdateReview обновляет review
func (h *UserInteractionHandler) UpdateReview(c *gin.Context) {
	ctx := c.Request.Context()

	// парисим id из ссылки
	reviewIDParam := c.Param("id")
	reviewID, err := strconv.Atoi(reviewIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем правильность запроса
	var req models.ReviewReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.service.UpdateReview(ctx, reviewID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// отправляем ответ
	c.JSON(http.StatusCreated, gin.H{"message": "review updated"})
}

// DeleteReview удаляет review
func (h *UserInteractionHandler) DeleteReview(c *gin.Context) {
	ctx := c.Request.Context()

	// получаем id
	reviewIDParam := c.Param("id")
	reviewID, err := strconv.Atoi(reviewIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// удаляем из таблицы
	if err = h.service.DeleteReview(ctx, reviewID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review deleted"})
}

// RegistrationOnEvent регистрация на ивент
func (h *UserInteractionHandler) RegistrationOnEvent(c *gin.Context) {

	// получаем id
	eventIDParam := c.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// берем userID и userName из Header
	if err = h.service.RegistrationOnEvent(c, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registration on"})
}

// DeleteRegistration удаление регистрации
func (h *UserInteractionHandler) DeleteRegistration(c *gin.Context) {
	// получаем id
	eventIDParam := c.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// берем userID из Header
	if err = h.service.DeleteRegistration(c, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "delete registration on"})
}

// GetRegistration получаем зарегистрированных пользователей на event
func (h *UserInteractionHandler) GetRegistration(c *gin.Context) {
	ctx := c.Request.Context()

	// получаем id
	eventIDParam := c.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// и собираем все registrations
	var registrations []models.ParticipantResp

	if err = h.service.GetRegistrations(ctx, eventID, &registrations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, registrations)
}

// RegisterRoutes собираем все хендлеры в одну функцию
func (UI *UserInteractionHandler) RegisterRoutes() {
	// user-interact
	userInteractionRev := UI.router.Group("/user-interact")

	// делаем middleware для проверки регистрации пользователя
	userInteractionRev.Use(jwt.AuthMiddleware())

	userInteractionRev.POST("/", UI.CreateNewReviews)
	userInteractionRev.GET("/event/:id", UI.GetCurrentReviewsByEventID)
	userInteractionRev.PUT("/:id", UI.UpdateReview)
	userInteractionRev.DELETE("/:id", UI.DeleteReview)

	// registation
	userInteractionRegister := UI.router.Group("/registration")

	userInteractionRegister.Use(jwt.AuthMiddleware())

	userInteractionRegister.POST("/event/:id", UI.RegistrationOnEvent)
	userInteractionRegister.DELETE("/event/:id", UI.DeleteRegistration)
	userInteractionRegister.GET("/event/:id", UI.GetRegistration)
}
