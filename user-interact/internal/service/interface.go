package service

import (
	"context"
	"eventify/user-interact/internal/models"
)

type Service interface {
	CreateNewReviews(ctx context.Context, req models.ReviewReq) error
	GetCurrentReviewsByEventID(ctx context.Context, eventId int, req *[]models.ReviewResp) error
	UpdateReview(ctx context.Context, reviewID int, req models.ReviewReq) error
	DeleteReview(ctx context.Context, reviewID int) error
	RegistrationOnEvent(ctx context.Context, eventID int, userID int, username string) error
	DeleteRegistration(ctx context.Context, eventID int, userID int) error
	GetRegistrations(ctx context.Context, eventID int, registrations *[]models.ParticipantResp) error
}
