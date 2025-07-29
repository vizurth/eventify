package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
	"vizurth/eventify/user-interaction/internal/models"
)

type UserInteractionRepository struct {
	db *pgxpool.Pool
}

func NewUserInteractionRepository(db *pgxpool.Pool) *UserInteractionRepository {
	return &UserInteractionRepository{
		db: db,
	}
}

func (r *UserInteractionRepository) CreateNewReviews(ctx context.Context, req models.ReviewReq) error {
	_, err := r.db.Exec(ctx, `
	INSERT INTO schema_name.reviews (event_id, user_id, username, rating, comment, updated_at)
	VALUES ($1, $2, $3, $4, $5, NULL)`,
		req.EventID,
		req.UserID,
		req.Username,
		req.Rating,
		req.Comment,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserInteractionRepository) GetCurrentReviewsByEventID(ctx context.Context, eventId int, reviews *[]models.ReviewResp) error {
	rows, err := r.db.Query(ctx, `
	SELECT * FROM schema_name.reviews
	WHERE event_id = $1`, eventId)
	if err != nil {
		return err
	}

	for rows.Next() {
		var temp models.ReviewResp
		err = rows.Scan(
			&temp.ID,
			&temp.EventID,
			&temp.UserID,
			&temp.Username,
			&temp.Rating,
			&temp.Comment,
			&temp.CreatedAt,
			&temp.UpdatedAt)
		if err != nil {
			return err
		}
		*reviews = append(*reviews, temp)
	}

	return nil
}

func (r *UserInteractionRepository) UpdateReview(ctx context.Context, reviewID int, req models.ReviewReq) error {
	// обновляем базу
	_, err := r.db.Exec(ctx, `
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
		return err
	}

	return nil
}

func (r *UserInteractionRepository) DeleteReview(ctx context.Context, reviewID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM schema_name.reviews WHERE id = $1`, reviewID)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserInteractionRepository) RegistrationOnEvent(ctx context.Context, eventID, userID int, username string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO schema_name.event_participants(event_id, user_id, username) VALUES ($1, $2, $3)`, eventID, userID, username)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserInteractionRepository) DeleteRegistration(ctx context.Context, eventID, userID int) error {
	_, err := r.db.Exec(ctx, `
	DELETE FROM schema_name.event_participants
	WHERE user_id = $1 AND event_id = $2`, userID, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserInteractionRepository) GetRegistrations(ctx context.Context, eventID int, registrations *[]models.ParticipantResp) error {
	rows, err := r.db.Query(ctx, `SELECT * FROM schema_name.event_participants WHERE event_id = $1`, eventID)
	defer rows.Close()

	if err != nil {
		return err
	}

	// добавляем все в наш массив
	for rows.Next() {
		var temp models.ParticipantResp
		err = rows.Scan(&temp.ID, &temp.Username, &temp.EventID)
		if err != nil {
			return err
		}
		*registrations = append(*registrations, temp)
	}
}
