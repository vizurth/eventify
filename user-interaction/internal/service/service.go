package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"eventify/common/kafka"
	"eventify/user-interaction/internal/models"
	"eventify/user-interaction/internal/repository"
)

type UserInteractionService struct {
	repo     *repository.UserInteractionRepository
	producer *kafka.Producer
}

func NewUserInteractionService(repo *repository.UserInteractionRepository, producer *kafka.Producer) *UserInteractionService {
	return &UserInteractionService{
		repo:     repo,
		producer: producer,
	}
}

func (s *UserInteractionService) CreateNewReviews(ctx context.Context, req models.ReviewReq) error {
	if err := s.repo.CreateNewReviews(ctx, req); err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":   req.UserID,
		"user_name": req.Username,
		"event_id":  req.EventID,
		"rating":    req.Rating,
		"comment":   req.Comment,
	})

	return s.producer.SendMessage(ctx, "review.created", payload)
}

func (s *UserInteractionService) GetCurrentReviewsByEventID(ctx context.Context, eventId int, req *[]models.ReviewResp) error {
	return s.repo.GetCurrentReviewsByEventID(ctx, eventId, req)
}

func (s *UserInteractionService) UpdateReview(ctx context.Context, reviewID int, req models.ReviewReq) error {
	if err := s.repo.UpdateReview(ctx, reviewID, req); err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":  req.UserID,
		"event_id": req.EventID,
		"rating":   req.Rating,
		"comment":  req.Comment,
	})

	return s.producer.SendMessage(ctx, "review.updated", payload)
}

func (s *UserInteractionService) DeleteReview(ctx context.Context, reviewID int) error {
	// можно расширить, если нужно знать ID пользователя/события
	payload, _ := json.Marshal(map[string]interface{}{
		"review_id": reviewID,
	})

	if err := s.repo.DeleteReview(ctx, reviewID); err != nil {
		return err
	}

	return s.producer.SendMessage(ctx, "review.deleted", payload)
}

func (s *UserInteractionService) RegistrationOnEvent(c *gin.Context, eventId int) error {
	ctx := c.Request.Context()

	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	userIDint := userID.(int)
	usernameStr := username.(string)

	if err := s.repo.RegistrationOnEvent(ctx, eventId, userIDint, usernameStr); err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":   userIDint,
		"event_id":  eventId,
		"user_name": usernameStr,
	})

	return s.producer.SendMessage(ctx, "registration.created", payload)
}

func (s *UserInteractionService) DeleteRegistration(c *gin.Context, eventId int) error {
	ctx := c.Request.Context()
	userID, _ := c.Get("userID")
	userIDint := userID.(int)

	if err := s.repo.DeleteRegistration(ctx, eventId, userIDint); err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":  userIDint,
		"event_id": eventId,
	})

	return s.producer.SendMessage(ctx, "registration.deleted", payload)
}

func (s *UserInteractionService) GetRegistrations(ctx context.Context, eventID int, registrations *[]models.ParticipantResp) error {
	return s.repo.GetRegistrations(ctx, eventID, registrations)
}
