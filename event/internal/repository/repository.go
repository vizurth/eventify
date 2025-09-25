package repository

import (
	"context"
	"encoding/json"
	"eventify/common/logger"
	"eventify/event/internal/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type EventRepository struct {
	db    *pgxpool.Pool
	psql  sq.StatementBuilderType
	redis *redis.Client
}

func NewEventRepository(db *pgxpool.Pool, redis *redis.Client) *EventRepository {
	return &EventRepository{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar), redis: redis}
}

func (r *EventRepository) CreateEvent(ctx context.Context, req models.EventReq) error {
	//добавляем данные в бд
	query, args, err := r.psql.Insert("events").
		Columns("title", "description", "category", "city", "venue", "address", "start_time", "end_time", "organizer_id", "organizer_name", "organizer_email", "status").
		Values(req.Title,
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
			req.Status).
		ToSql()

	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	var lastEventID uint
	query, args, err = r.psql.Select("id").From("events").OrderBy("id DESC LIMIT 1").ToSql()

	err = r.db.QueryRow(ctx, query, args...).Scan(&lastEventID)
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	// цикл для добавления участников в базу данных
	for _, participant := range req.Participants {
		query, args, err = r.psql.Insert("event_participants").
			Columns("event_id", "user_id", "username").
			Values(lastEventID, participant.ID, participant.Username).ToSql()

		if err != nil {
			return fmt.Errorf("create event: create event_participants: %w", err)
		}
		_, err = r.db.Exec(ctx, query, args...)

		if err != nil {
			return fmt.Errorf("create event: create event_participants: %w", err)
		}
	}

	return nil
}

func (r *EventRepository) GetEvents(ctx context.Context, events *[]models.EventResp) error {
	log := logger.GetOrCreateLoggerFromCtx(ctx)
	// Redis
	cacheKey := "events:list"

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), events); err == nil {
			log.Info(ctx, "list found from redis")
			return nil
		}
	}

	// Postgres
	query, args, err := r.psql.
		Select(
			"id", "title", "description", "category",
			"city", "venue", "address",
			"start_time", "end_time",
			"organizer_id", "organizer_name", "organizer_email",
			"status", "created_at").
		From("events").
		ToSql()
	if err != nil {
		return fmt.Errorf("get events: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("get events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e models.EventResp
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.Category,
			&e.Location.City, &e.Location.Venue, &e.Location.Address,
			&e.StartTime, &e.EndTime,
			&e.Organizer.ID, &e.Organizer.Username, &e.Organizer.Email,
			&e.Status, &e.CreatedAt,
		); err != nil {
			return fmt.Errorf("get events: %w", err)
		}

		// Получаем участников через squirrel
		pQuery, pArgs, err := r.psql.
			Select("user_id", "username").
			From("event_participants").
			Where(sq.Eq{"event_id": e.ID}).
			ToSql()
		if err != nil {
			return fmt.Errorf("get events: %w", err)
		}

		pRows, err := r.db.Query(ctx, pQuery, pArgs...)
		if err != nil {
			return fmt.Errorf("get events: %w", err)
		}
		var participants []models.Participant
		for pRows.Next() {
			var p models.Participant
			if err := pRows.Scan(&p.ID, &p.Username); err != nil {
				pRows.Close()
				return fmt.Errorf("get events: %w", err)
			}
			participants = append(participants, p)
		}
		pRows.Close()

		e.Participants = participants
		*events = append(*events, e)
	}

	data, _ := json.Marshal(*events)
	_ = r.redis.Set(ctx, cacheKey, data, 30*time.Second)
	log.Info(ctx, "list event set to redis")

	return nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, eventID int, e *models.EventResp) error {
	log := logger.GetOrCreateLoggerFromCtx(ctx)
	cachedEvent, err := r.redis.Get(ctx, "event:"+strconv.Itoa(eventID)).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedEvent), e); err == nil {
			log.Info(ctx, "Event found from redis", zap.Int("id", eventID))
			return nil
		}
	}

	// Событие по ID через squirrel
	query, args, err := r.psql.
		Select(
			"id", "title", "description", "category",
			"city", "venue", "address",
			"start_time", "end_time",
			"organizer_id", "organizer_name", "organizer_email",
			"status", "created_at").
		From("events").
		Where(sq.Eq{"id": eventID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("get event by id: %w", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&e.ID, &e.Title, &e.Description, &e.Category,
		&e.Location.City, &e.Location.Venue, &e.Location.Address,
		&e.StartTime, &e.EndTime,
		&e.Organizer.ID, &e.Organizer.Username, &e.Organizer.Email,
		&e.Status, &e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("get event by id: %w", err)
	}

	// Участники через squirrel
	pQuery, pArgs, err := r.psql.
		Select("user_id", "username").
		From("event_participants").
		Where(sq.Eq{"event_id": e.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("get event by id: %w", err)
	}

	rows, err := r.db.Query(ctx, pQuery, pArgs...)
	if err != nil {
		return fmt.Errorf("get event by id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Participant
		if err := rows.Scan(&p.ID, &p.Username); err != nil {
			return fmt.Errorf("get event by id: %w", err)
		}
		e.Participants = append(e.Participants, p)
	}

	data, _ := json.Marshal(*e)
	_ = r.redis.Set(ctx, "event:"+strconv.Itoa(eventID), data, 15*time.Minute)
	log.Info(ctx, "Event set to redis", zap.Int("id", eventID))

	return nil
}

func (r *EventRepository) CheckUserRegistration(ctx context.Context, eventID, userID int, e *models.EventResp) error {
	query, args, err := r.psql.
		Select("EXISTS(SELECT 1 FROM event_participants WHERE event_id = ? AND user_id = ?)").
		Where(sq.Eq{"event_id": eventID, "user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("check user registration: %w", err)
	}

	var exists bool
	if err = r.db.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		return fmt.Errorf("check user registration: %w", err)
	}

	e.IsRegistered = exists
	return nil
}
