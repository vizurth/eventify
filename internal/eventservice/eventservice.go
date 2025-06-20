package eventservice

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vizurth/eventify/models"
	"net/http"
	"strconv"
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
	var req models.EventReq
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
	var lastEventID uint
	query := `SELECT id FROM schema_name.events ORDER BY id DESC LIMIT 1`

	err = es.db.QueryRow(ctx, query).Scan(&lastEventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get last event ID: " + err.Error()})
		return
	}
	// шаблон для добавления юсера в схему
	insertParticipantQuery := `
        INSERT INTO schema_name.event_participants (event_id, user_id, username)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;
    `

	// цикл для добавления участников в базу данных
	for _, participant := range req.Participants {
		_, err := es.db.Exec(ctx, insertParticipantQuery,
			lastEventID,
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

// GetEvents получаем все ивенты
func (es *EventService) GetEvents(c *gin.Context) {
	ctx := c.Request.Context()
	// получаем все events
	var events []models.EventResp
	rows, err := es.db.Query(ctx, `
	SELECT * FROM schema_name.events
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer rows.Close()
	// проходимся по всем строкам
	for rows.Next() {
		var e models.EventResp
		err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.Category,
			&e.Location.City,
			&e.Location.Venue,
			&e.Location.Address,
			&e.StartTime,
			&e.EndTime,
			&e.Organizer.ID,
			&e.Organizer.Username,
			&e.Organizer.Email,
			&e.Status,
			&e.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// запрашиваем участников event
		participantsRows, err := es.db.Query(ctx, `
        SELECT user_id, username 
        FROM schema_name.event_participants
        WHERE event_id = $1
    `, e.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get participants: " + err.Error()})
			return
		}
		// записываем в нашу переменную
		var participants []models.Participant
		for participantsRows.Next() {
			var p models.Participant
			err := participantsRows.Scan(&p.ID, &p.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			participants = append(participants, p)
		}
		participantsRows.Close()
		// добавляем в e
		e.Participants = participants
		// заполняем events
		events = append(events, e)
	}

	c.JSON(http.StatusOK, events)
}

// GetEventById получение ивента по id
func (es *EventService) GetEventById(c *gin.Context) {
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
	err = es.db.QueryRow(ctx, `
        SELECT 
            id, title, description, category, city, venue, address,
            start_time, end_time, organizer_id, organizer_name, organizer_email,
            status, created_at
        FROM schema_name.events
        WHERE id = $1
    `, eventID).Scan(
		&e.ID,
		&e.Title, &e.Description, &e.Category,
		&e.Location.City, &e.Location.Venue, &e.Location.Address,
		&e.StartTime, &e.EndTime,
		&e.Organizer.ID, &e.Organizer.Username, &e.Organizer.Email,
		&e.Status, &e.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// считываем участников и записываем в event
	participants, err := es.db.Query(ctx, `
        SELECT user_id, username FROM schema_name.event_participants
        WHERE event_id = $1
    `, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get participants: " + err.Error()})
		return
	}
	defer participants.Close()
	for participants.Next() {
		var p models.Participant
		if err := participants.Scan(&p.ID, &p.Username); err == nil {
			e.Participants = append(e.Participants, p)
		}
	}

	c.JSON(http.StatusOK, e)

}

// RegisterRoutes собираем все хендлеры в одну функцию
func (es *EventService) RegisterRoutes() {
	events := es.router.Group("/events")
	events.POST("/", es.CreateEvent)
	events.GET("/", es.GetEvents)
	events.GET("/:id", es.GetEventById)
}
