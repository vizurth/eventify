package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"vizurth/eventify/models"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) CreateEvent(ctx context.Context, req models.EventReq) error {
	//добавляем данные в бд
	_, err := r.db.Exec(ctx,
		`INSERT INTO schema_name.event(title, description, category, city, venue, address, start_time, end_time, organizer_id, organizer_name, organizer_email, status)
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
		return err
	}
	var lastEventID uint
	query := `SELECT id FROM schema_name.event ORDER BY id DESC LIMIT 1`

	err = r.db.QueryRow(ctx, query).Scan(&lastEventID)
	if err != nil {
		return err
	}
	// шаблон для добавления юсера в схему
	insertParticipantQuery := `
      INSERT INTO schema_name.event_participants (event_id, user_id, username)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;
  `

	// цикл для добавления участников в базу данных
	for _, participant := range req.Participants {
		_, err = r.db.Exec(ctx, insertParticipantQuery,
			lastEventID,
			participant.ID,
			participant.Username,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *EventRepository) GetEvents(ctx context.Context, events *[]models.EventResp) error {
	rows, err := r.db.Query(ctx, `
	SELECT * FROM schema_name.event
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	// проходимся по всем строкам
	for rows.Next() {
		var e models.EventResp
		err = rows.Scan(
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
			return err
		}

		// запрашиваем участников event
		participantsRows, err := r.db.Query(ctx, `
        SELECT user_id, username 
        FROM schema_name.event_participants
        WHERE event_id = $1
    `, e.ID)
		if err != nil {
			return err
		}
		// записываем в нашу переменную
		var participants []models.Participant
		for participantsRows.Next() {
			var p models.Participant
			err = participantsRows.Scan(&p.ID, &p.Username)
			if err != nil {
				return err
			}
			participants = append(participants, p)
		}
		participantsRows.Close()
		// добавляем в e
		e.Participants = participants
		// заполняем event
		*events = append(*events, e)
	}
	return nil
}

func (r *EventRepository) GetEventByID(ctx context.Context, eventID int, e *models.EventResp) error {
	err := r.db.QueryRow(ctx, `
	       SELECT
	           id, title, description, category, city, venue, address,
	           start_time, end_time, organizer_id, organizer_name, organizer_email,
	           status, created_at
	       FROM schema_name.event
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
		return err
	}
	// считываем участников и записываем в event
	participants, err := r.db.Query(ctx, `
	       SELECT user_id, username FROM schema_name.event_participants
	       WHERE event_id = $1
	   `, eventID)
	if err != nil {
		return err
	}
	defer participants.Close()
	for participants.Next() {
		var p models.Participant
		if err := participants.Scan(&p.ID, &p.Username); err == nil {
			e.Participants = append(e.Participants, p)
		}
	}
}
