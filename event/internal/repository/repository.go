package repository

import (
	"context"
	"eventify/event/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
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
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	if err != nil {
		return err
	}
	var lastEventID uint
	query, args, err = r.psql.Select("id").From("events").OrderBy("id DESC LIMIT 1").ToSql()

	err = r.db.QueryRow(ctx, query, args...).Scan(&lastEventID)
	if err != nil {
		return err
	}

	// цикл для добавления участников в базу данных
	for _, participant := range req.Participants {
		query, args, err = r.psql.Insert("event_participants").
			Columns("event_id", "user_id", "username").
			Values(lastEventID, participant.ID, participant.Username).ToSql()

		if err != nil {
			return err
		}
		_, err = r.db.Exec(ctx, query, args...)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *EventRepository) GetEvents(ctx context.Context, events *[]models.EventResp) error {
	// Собираем запрос через squirrel
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
		return err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return err
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
			return err
		}

		// Получаем участников через squirrel
		pQuery, pArgs, err := r.psql.
			Select("user_id", "username").
			From("event_participants").
			Where(sq.Eq{"event_id": e.ID}).
			ToSql()
		if err != nil {
			return err
		}

		pRows, err := r.db.Query(ctx, pQuery, pArgs...)
		if err != nil {
			return err
		}
		var participants []models.Participant
		for pRows.Next() {
			var p models.Participant
			if err := pRows.Scan(&p.ID, &p.Username); err != nil {
				pRows.Close()
				return err
			}
			participants = append(participants, p)
		}
		pRows.Close()

		e.Participants = participants
		*events = append(*events, e)
	}

	return nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, eventID int, e *models.EventResp) error {
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
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&e.ID, &e.Title, &e.Description, &e.Category,
		&e.Location.City, &e.Location.Venue, &e.Location.Address,
		&e.StartTime, &e.EndTime,
		&e.Organizer.ID, &e.Organizer.Username, &e.Organizer.Email,
		&e.Status, &e.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Участники через squirrel
	pQuery, pArgs, err := r.psql.
		Select("user_id", "username").
		From("event_participants").
		Where(sq.Eq{"event_id": e.ID}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Query(ctx, pQuery, pArgs...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Participant
		if err := rows.Scan(&p.ID, &p.Username); err != nil {
			return err
		}
		e.Participants = append(e.Participants, p)
	}

	return nil
}

func (r *EventRepository) CheckUserRegistration(ctx context.Context, eventID, userID int, e *models.EventResp) error {
	query, args, err := r.psql.
		Select("EXISTS(SELECT 1 FROM event_participants WHERE event_id = ? AND user_id = ?)").
		Where(sq.Eq{"event_id": eventID, "user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	var exists bool
	if err = r.db.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		return err
	}

	e.IsRegistered = exists
	return nil
}
