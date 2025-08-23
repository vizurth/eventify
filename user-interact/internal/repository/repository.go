package repository

import (
	"context"
	"eventify/user-interact/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserInteractionRepository struct {
	db   *pgxpool.Pool
	pgql sq.StatementBuilderType
}

func NewUserInteractionRepository(db *pgxpool.Pool) *UserInteractionRepository {
	return &UserInteractionRepository{
		db:   db,
		pgql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserInteractionRepository) CreateNewReviews(ctx context.Context, req models.ReviewReq) error {
	sql, args, err := r.pgql.Insert("reviews").
		Columns("event_id", "user_id", "username", "rating", "comment", "updated_at").
		Values(req.EventID, req.UserID, req.Username, req.Rating, req.Comment, nil).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *UserInteractionRepository) GetCurrentReviewsByEventID(ctx context.Context, eventId int, reviews *[]models.ReviewResp) error {
	sql, args, err := r.pgql.Select("*").
		From("reviews").
		Where(sq.Eq{"event_id": eventId}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var temp models.ReviewResp
		if err := rows.Scan(
			&temp.ID,
			&temp.EventID,
			&temp.UserID,
			&temp.Username,
			&temp.Rating,
			&temp.Comment,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		); err != nil {
			return err
		}
		*reviews = append(*reviews, temp)
	}
	return nil
}

func (r *UserInteractionRepository) UpdateReview(ctx context.Context, reviewID int, req models.ReviewReq) error {
	sql, args, err := r.pgql.Update("reviews").
		Set("rating", req.Rating).
		Set("comment", req.Comment).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": reviewID, "username": req.Username, "user_id": req.UserID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *UserInteractionRepository) DeleteReview(ctx context.Context, reviewID int) error {
	sql, args, err := r.pgql.Delete("reviews").
		Where(sq.Eq{"id": reviewID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *UserInteractionRepository) RegistrationOnEvent(ctx context.Context, eventID, userID int, username string) error {
	sql, args, err := r.pgql.Insert("event_participants").
		Columns("event_id", "user_id", "username").
		Values(eventID, userID, username).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *UserInteractionRepository) DeleteRegistration(ctx context.Context, eventID, userID int) error {
	sql, args, err := r.pgql.Delete("event_participants").
		Where(sq.Eq{"event_id": eventID, "user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *UserInteractionRepository) GetRegistrations(ctx context.Context, eventID int, registrations *[]models.ParticipantResp) error {
	sql, args, err := r.pgql.Select("*").
		From("event_participants").
		Where(sq.Eq{"event_id": eventID}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var temp models.ParticipantResp
		if err := rows.Scan(&temp.ID, &temp.Username, &temp.EventID); err != nil {
			return err
		}
		*registrations = append(*registrations, temp)
	}

	return nil
}
