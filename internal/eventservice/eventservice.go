package eventservice

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vizurth/eventify/models"
	"net/http"
)

type EventService struct {
	db     *pgxpool.Pool
	router *gin.Engine
}

func NewEventService(db *pgxpool.Pool, router *gin.Engine) *EventService {
	return &EventService{
		db:     db,
		router: router,
	}
}

// CreateEvent создает event по запросу /events/
func (es *EventService) CreateEvent(c *gin.Context) {
	// прокидываем ctx для работы с базой данных
	ctx := c.Request.Context()

	// проверяем корректность JSON запроса
	var req models.Event
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// добавляем данные в бд
	_, err := es.db.Exec(ctx,
		`INSERT INTO schema_name.events(title, description, category, city, venue, address, start_time, end_time, organizer_id, organizer_name, organizer_email, status)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		req.Title,
		req.Description,
		req.Category,
		req.Location.City,
		req.Location.Venue,
		req.Location.Address,
		req.StartTime, // Время начала события (в top‑level, а не в Location)
		req.EndTime,   // Время окончания события
		req.Organizer.ID,
		req.Organizer.Username,
		req.Organizer.Email,
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// шаблон для добавления юсера в схему
	insertParticipantQuery := `
        INSERT INTO schema_name.event_participants
            (event_id, user_id, username)
        VALUES
            ($1, $2, $3)
    `

	// цикл для добавления участников в базу данных
	for _, participant := range req.Participants {
		_, err := es.db.Exec(ctx, insertParticipantQuery,
			participant.ID,
			participant.Username,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert participant: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created"})
}

// RegisterRoutes собираем все хендлеры в одну функцию
func (es *EventService) RegisterRoutes() {
	events := es.router.Group("/events")
	events.POST("/", es.CreateEvent)
}
