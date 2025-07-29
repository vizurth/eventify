package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"vizurth/eventify/user-interaction/internal/models"
	"vizurth/eventify/user-interaction/internal/repository"
)

type UserInteractionService struct {
	repo *repository.UserInteractionRepository
}

func NewUserInteractionService(repo *repository.UserInteractionRepository) *UserInteractionService {
	return &UserInteractionService{
		repo: repo,
	}
}

func (s *UserInteractionService) CreateNewReviews(ctx context.Context, req models.ReviewReq) error {
	if err := s.repo.CreateNewReviews(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *UserInteractionService) GetCurrentReviewsByEventID(ctx context.Context, eventId int, req *[]models.ReviewResp) error {
	if err := s.repo.GetCurrentReviewsByEventID(ctx, eventId, req); err != nil {
		return err
	}
	return nil
}

func (s *UserInteractionService) UpdateReview(ctx context.Context, reviewID int, req models.ReviewReq) error {
	if err := s.repo.UpdateReview(ctx, reviewID, req); err != nil {
		return err
	}

	return nil
}

func (s *UserInteractionService) DeleteReview(ctx context.Context, reviewID int) error {
	if err := s.repo.DeleteReview(ctx, reviewID); err != nil {
		return err
	}

	return nil
}

func (s *UserInteractionService) RegistrationOnEvent(c *gin.Context, eventId int) error {
	ctx := c.Request.Context()

	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	userIDint := userID.(int)
	usernameStr := username.(string)

	// добавляем в таблицу user
	if err := s.repo.RegistrationOnEvent(ctx, eventId, userIDint, usernameStr); err != nil {
		return err
	}
	return nil
}

func (s *UserInteractionService) DeleteRegistration(c *gin.Context, eventId int) error {
	ctx := c.Request.Context()
	userID, _ := c.Get("userID")

	userIDint := userID.(int)

	// удаляем из таблицы
	if err := s.repo.DeleteRegistration(ctx, eventId, userIDint); err != nil {
		return err
	}
	return nil
}

func (s *UserInteractionService) GetRegistrations(ctx context.Context, eventID int, registrations *[]models.ParticipantResp) error {
	if err := s.repo.GetRegistrations(ctx, eventID, registrations); err != nil {
		return err
	}
	return nil
}
