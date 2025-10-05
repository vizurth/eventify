package repository

import (
	"context"
	"eventify/user-interact/internal/models"
)

type Repository interface {
	CreateNewReviews(ctx context.Context, req models.ReviewReq) error
	GetCurrentReviewsByEventID(ctx context.Context, eventId int, reviews *[]models.ReviewResp) error
	UpdateReview(ctx context.Context, reviewID int, req models.ReviewReq) error
	DeleteReview(ctx context.Context, reviewID int) error
	RegistrationOnEvent(ctx context.Context, eventID, userID int, username string) error
	DeleteRegistration(ctx context.Context, eventID, userID int) error
	GetRegistrations(ctx context.Context, eventID int, registrations *[]models.ParticipantResp) error
}
