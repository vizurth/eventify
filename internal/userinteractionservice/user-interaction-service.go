package userInteractionService

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vizurth/eventify/internal/middleware"
	"github.com/vizurth/eventify/models"
	"net/http"
	"strconv"
	"time"
)

type UserInteractionService struct {
	db     *pgxpool.Pool
	router *gin.Engine
	secret []byte
}

func NewUserInteractionService(db *pgxpool.Pool, router *gin.Engine, secret []byte) *UserInteractionService {
	return &UserInteractionService{
		db:     db,
		router: router,
		secret: secret,
	}
}

// CreateNewReviews создаем новые отзывы
func (UI *UserInteractionService) CreateNewReviews(c *gin.Context) {
	ctx := c.Request.Context()

	// проверяем корректность запроса
	var req models.ReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// добавляем review в таблицу
	_, err := UI.db.Exec(ctx, `
	INSERT INTO schema_name.reviews (event_id, user_id, username, rating, comment, updated_at)
	VALUES ($1, $2, $3, $4, $5, NULL)`,
		req.EventID,
		req.UserID,
		req.Username,
		req.Rating,
		req.Comment,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review created"})
}

// GetCurrentReviewsByEventID получение review по корректному id
func (UI *UserInteractionService) GetCurrentReviewsByEventID(c *gin.Context) {
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

	rows, err := UI.db.Query(ctx, `
	SELECT * FROM schema_name.reviews
	WHERE event_id = $1`, eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	for rows.Next() {
		var r models.ReviewResp
		err = rows.Scan(
			&r.ID,
			&r.EventID,
			&r.UserID,
			&r.Username,
			&r.Rating,
			&r.Comment,
			&r.CreatedAt,
			&r.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reviews = append(reviews, r)
	}

	c.JSON(http.StatusOK, reviews)
}

// UpdateReview обновляет review
func (UI *UserInteractionService) UpdateReview(c *gin.Context) {
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

	// обновляем базу
	_, err = UI.db.Exec(ctx, `
	UPDATE schema_name.reviews
	SET rating = $1, comment = $2, updated_at = $3
	WHERE id = $4 AND username = $5 AND user_id = $6;`,
		req.Rating,
		req.Comment,
		time.Now(),
		reviewID,
		req.Username,
		req.UserID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// отправляем ответ
	c.JSON(http.StatusCreated, gin.H{"message": "review updated"})
}

// DeleteReview удаляет review
func (UI *UserInteractionService) DeleteReview(c *gin.Context) {
	ctx := c.Request.Context()

	// получаем id
	reviewIDParam := c.Param("id")
	reviewID, err := strconv.Atoi(reviewIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// удаляем из таблицы
	_, err = UI.db.Exec(ctx, `DELETE FROM schema_name.reviews WHERE id = $1`, reviewID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review deleted"})
}

// RegistrationOnEvent регистрация на ивент
func (UI *UserInteractionService) RegistrationOnEvent(c *gin.Context) {
	ctx := c.Request.Context()

	// получаем id
	eventIDParam := c.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// берем userID и userName из Header
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	userIDint := userID.(int)
	usernamestr := username.(string)

	// добавляем в таблицу user
	_, err = UI.db.Exec(ctx, `INSERT INTO schema_name.event_participants(event_id, user_id, username) VALUES ($1, $2, $3)`, eventID, userIDint, usernamestr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registration on"})
}

// DeleteRegistration удаление регистрации
func (UI *UserInteractionService) DeleteRegistration(c *gin.Context) {
	ctx := c.Request.Context()

	// получаем id
	eventIDParam := c.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// берем userID из Header
	userID, _ := c.Get("userID")

	userIDint := userID.(int)

	// удаляем из таблицы
	_, err = UI.db.Exec(ctx, `
	DELETE FROM schema_name.event_participants 
	WHERE user_id = $1 AND event_id = $2`, userIDint, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "delete registration on"})
}

// GetRegistration получаем зарегистрированных пользователей на event
func (UI *UserInteractionService) GetRegistration(c *gin.Context) {
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

	rows, err := UI.db.Query(ctx, `SELECT * FROM schema_name.event_participants WHERE event_id = $1`, eventID)
	defer rows.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// добавляем все в наш массив
	for rows.Next() {
		var r models.ParticipantResp
		err = rows.Scan(&r.ID, &r.Username, &r.EventID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		registrations = append(registrations, r)
	}

	c.JSON(http.StatusOK, registrations)
}

// RegisterRoutes собираем все хендлеры в одну функцию
func (UI *UserInteractionService) RegisterRoutes() {
	// reviews
	userInteraction := UI.router.Group("/")

	// делаем middleware для проверки регистрации пользователя
	userInteraction.Use(middleware.AuthMiddleWareDefault(UI.secret))

	userInteraction.POST("/reviews", UI.CreateNewReviews)
	userInteraction.GET("/reviews/event/:id", UI.GetCurrentReviewsByEventID)
	userInteraction.PUT("/reviews/:id", UI.UpdateReview)
	userInteraction.DELETE("/reviews/:id", UI.DeleteReview)

	//registation
	userInteraction.POST("/registration/event/:id", UI.RegistrationOnEvent)
	userInteraction.DELETE("/registration/event/:id", UI.DeleteRegistration)
	userInteraction.GET("/registration/event/:id", UI.GetRegistration)
}
