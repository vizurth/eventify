package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"eventify/event/internal/models"
	"eventify/event/internal/service"
)

type EventHandler struct {
	service *service.EventService
	router  *gin.Engine
}

func NewEventHandler(service *service.EventService, router *gin.Engine) *EventHandler {
	return &EventHandler{
		service: service,
		router:  gin.Default(),
	}
}

// CreateEvent создает event по запросу /event/
func (h *EventHandler) CreateEvent(c *gin.Context) {
	// прокидываем ctx для работы с базой данных
	ctx := c.Request.Context()

	// проверяем корректность JSON запроса
	var req models.EventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateEvent(ctx, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created"})
}

// GetEvents получаем все ивенты
func (h *EventHandler) GetEvents(c *gin.Context) {
	ctx := c.Request.Context()

	// получаем все event
	var events []models.EventResp

	if err := h.service.GetEvents(ctx, &events); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEventById получение ивента по id
func (h *EventHandler) GetEventById(c *gin.Context) {
	ctx := c.Request.Context()

	// считываем eventID
	eventId := c.Param("id")
	eventID, err := strconv.Atoi(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}
	// запрашиваем event по ID
	var e models.EventResp

	if err = h.service.GetEventByID(ctx, eventID, &e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, e)

}

// RegisterRoutes собираем все хендлеры в одну функцию
func (h *EventHandler) RegisterRoutes() {
	events := h.router.Group("/")
	events.POST("/", h.CreateEvent)
	events.GET("/", h.GetEvents)
	events.GET("/:id", h.GetEventById)
}
