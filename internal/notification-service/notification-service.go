package notificationService

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vizurth/eventify/models"
	"net/http"
	"strconv"
)

type NotificationService struct {
	db     *pgxpool.Pool
	router *gin.Engine
	secret []byte
}

func NewNotificationService(db *pgxpool.Pool, router *gin.Engine, secret []byte) *NotificationService {
	return &NotificationService{
		db:     db,
		router: router,
		secret: secret,
	}
}

// SendNotification отправка новых уведомлений
func (n *NotificationService) SendNotification(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.NotificationToUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := n.db.Exec(ctx, `INSERT INTO schema_name.notifications (user_id, message, is_read) VALUES ($1, $2, $3)`, req.UserID, req.Message, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "notification sent"})
}

// GetUserNotification получение пользовательских уведомлений
func (n *NotificationService) GetUserNotification(c *gin.Context) {
	ctx := c.Request.Context()

	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := n.db.Query(ctx, `SELECT * FROM schema_name.notifications WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var notifications []models.NotificationGetResp

	for rows.Next() {
		var n models.NotificationGetResp

		rows.Scan(&n.NotificationID, &n.UserID, &n.Message, &n.IsRead, &n.CreatedAt)

		notifications = append(notifications, n)
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// ReadNotification читаем уведомление
func (n *NotificationService) ReadNotification(c *gin.Context) {
	ctx := c.Request.Context()
	notificationIDParam := c.Param("id")
	notificationID, err := strconv.Atoi(notificationIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = n.db.Exec(ctx, `
	UPDATE schema_name.notifications 
	SET is_read = $1 
	WHERE id = $2
	`, true, notificationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "notification read"})
}

// DeleteNotification удалить уведомление из таблицы
func (n *NotificationService) DeleteNotification(c *gin.Context) {
	ctx := c.Request.Context()

	notificationIDParam := c.Param("id")
	notificationID, err := strconv.Atoi(notificationIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = n.db.Exec(ctx, `DELETE FROM schema_name.notifications WHERE id = $1`, notificationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "notification delete"})
}

func (n *NotificationService) RegisterRoutes() {
	notifications := n.router.Group("/")
	//notifications.Use(middleware.AuthMiddleWareDefault(n.secret))

	notifications.POST("/send", n.SendNotification)
	notifications.GET("/user/:id", n.GetUserNotification)
	notifications.PUT("/:id/read", n.ReadNotification)
	notifications.DELETE("/:id", n.DeleteNotification)

	//subscriptions := n.router.Group("/subscriptions")

}
