package userInteractionService

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vizurth/eventify/models"
	"net/http"
	"strconv"
)

type UserInteractionService struct {
	db     *pgxpool.Pool
	router *gin.Engine
}

func NewUserInteractionService(db *pgxpool.Pool, router *gin.Engine) *UserInteractionService {
	return &UserInteractionService{
		db:     db,
		router: router,
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

func (UI *UserInteractionService) RegisterRoutes() {
	userInteraction := UI.router.Group("/")
	userInteraction.POST("/reviews", UI.CreateNewReviews)
	userInteraction.POST("/reviews/event/:id", UI.GetCurrentReviewsByEventID)
}
