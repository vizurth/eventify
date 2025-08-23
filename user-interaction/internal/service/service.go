package service

import (
	"context"
	"eventify/user-interaction/internal/models"
	"eventify/user-interaction/internal/repository"
)

type UserInteractionService struct {
	repo *repository.UserInteractionRepository
}

func NewUserInteractionService(repo *repository.UserInteractionRepository) *UserInteractionService {
	return &UserInteractionService{repo: repo}
}

func (s *UserInteractionService) CreateNewReviews(ctx context.Context, req models.ReviewReq) error {
	return s.repo.CreateNewReviews(ctx, req)
}

func (s *UserInteractionService) GetCurrentReviewsByEventID(ctx context.Context, eventId int, req *[]models.ReviewResp) error {
	return s.repo.GetCurrentReviewsByEventID(ctx, eventId, req)
}

func (s *UserInteractionService) UpdateReview(ctx context.Context, reviewID int, req models.ReviewReq) error {
	return s.repo.UpdateReview(ctx, reviewID, req)
}

func (s *UserInteractionService) DeleteReview(ctx context.Context, reviewID int) error {
	return s.repo.DeleteReview(ctx, reviewID)
}

func (s *UserInteractionService) RegistrationOnEvent(ctx context.Context, eventId int, userID int, username string) error {
	return s.repo.RegistrationOnEvent(ctx, eventId, userID, username)
}

func (s *UserInteractionService) DeleteRegistration(ctx context.Context, eventId int, userID int) error {
	return s.repo.DeleteRegistration(ctx, eventId, userID)
}

func (s *UserInteractionService) GetRegistrations(ctx context.Context, eventID int, registrations *[]models.ParticipantResp) error {
	return s.repo.GetRegistrations(ctx, eventID, registrations)
}
