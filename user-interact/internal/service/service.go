package service

import (
	"context"
	"encoding/json"
	mykafka "eventify/common/kafka"
	"eventify/common/logger"
	"eventify/common/retry"
	"eventify/user-interact/internal/models"
	"eventify/user-interact/internal/repository"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type UserInteractionService struct {
	repo *repository.UserInteractionRepository

	reviewCreatedW       *mykafka.Writer
	reviewUpdatedW       *mykafka.Writer
	reviewDeletedW       *mykafka.Writer
	registrationCreatedW *mykafka.Writer
	registrationDeleteW  *mykafka.Writer
}

func NewUserInteractionService(ctx context.Context, repo *repository.UserInteractionRepository, cfg mykafka.Config) *UserInteractionService {
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	topics := []string{
		"review-created",
		"review-updated",
		"review-deleted",
		"registration-created",
		"registration-deleted",
	}
	for _, topic := range topics {
		if err := mykafka.CreateTopicWithRetry(cfg, topic, 1, 1, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 2}); err != nil {
			log.Error(ctx, "failed to create topic", zap.String("topic", topic), zap.Error(err))
			return nil
		}
	}

	return &UserInteractionService{
		repo:                 repo,
		reviewCreatedW:       mykafka.NewWriter(ctx, cfg, "review-created"),
		reviewUpdatedW:       mykafka.NewWriter(ctx, cfg, "review-updated"),
		reviewDeletedW:       mykafka.NewWriter(ctx, cfg, "review-deleted"),
		registrationCreatedW: mykafka.NewWriter(ctx, cfg, "registration-created"),
		registrationDeleteW:  mykafka.NewWriter(ctx, cfg, "registration-deleted")}
}

func (s *UserInteractionService) CreateNewReviews(ctx context.Context, req models.ReviewReq) error {
	log := logger.GetOrCreateLoggerFromCtx(ctx)
	if err := s.repo.CreateNewReviews(ctx, req); err != nil {
		return err
	}

	eventPayload := map[string]string{
		"event_id": strconv.Itoa(int(req.EventID)),
		"username": req.Username,
		"rating":   strconv.Itoa(int(req.Rating)),
		"comment":  req.Comment,
	}

	value, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}
	err = s.reviewCreatedW.SendWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 3}, []byte("review.created"), value)

	if err != nil {
		log.Error(ctx, "failed to send review.created", zap.Error(err))
		return err
	}
	log.Info(ctx, "successfully review.created")

	return nil
}

func (s *UserInteractionService) GetCurrentReviewsByEventID(ctx context.Context, eventId int, req *[]models.ReviewResp) error {
	return s.repo.GetCurrentReviewsByEventID(ctx, eventId, req)
}

func (s *UserInteractionService) UpdateReview(ctx context.Context, reviewID int, req models.ReviewReq) error {
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	// обновляем в БД
	if err := s.repo.UpdateReview(ctx, reviewID, req); err != nil {
		return err
	}

	eventPayload := map[string]string{
		"review_id": strconv.Itoa(reviewID),
		"event_id":  strconv.Itoa(int(req.EventID)),
		"username":  req.Username,
		"rating":    strconv.Itoa(int(req.Rating)),
		"comment":   req.Comment,
	}

	value, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	if err := s.reviewUpdatedW.SendWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 3}, []byte("review.updated"), value); err != nil {
		log.Error(ctx, "failed to send review.updated", zap.Error(err))
		return err
	}

	log.Info(ctx, "successfully review.updated", zap.String("review_id", strconv.Itoa(reviewID)))
	return nil
}

func (s *UserInteractionService) DeleteReview(ctx context.Context, reviewID int) error {
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := s.repo.DeleteReview(ctx, reviewID); err != nil {
		return err
	}

	eventPayload := map[string]string{
		"review_id": strconv.Itoa(reviewID),
	}

	value, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	if err := s.reviewDeletedW.SendWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 3}, []byte("review.deleted"), value); err != nil {
		log.Error(ctx, "failed to send review.deleted", zap.Error(err))
		return err
	}

	log.Info(ctx, "successfully review.deleted", zap.String("review_id", strconv.Itoa(reviewID)))
	return nil
}

func (s *UserInteractionService) RegistrationOnEvent(ctx context.Context, eventID int, userID int, username string) error {
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := s.repo.RegistrationOnEvent(ctx, eventID, userID, username); err != nil {
		return err
	}

	eventPayload := map[string]string{
		"event_id": strconv.Itoa(eventID),
		"user_id":  strconv.Itoa(userID),
		"username": username,
	}

	value, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	if err := s.registrationCreatedW.SendWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 3}, []byte("registration.created"), value); err != nil {
		log.Error(ctx, "failed to send registration.created", zap.Error(err))
		return err
	}

	log.Info(ctx, "successfully registration.created", zap.String("event_id", strconv.Itoa(eventID)), zap.String("user_id", strconv.Itoa(userID)))
	return nil
}

func (s *UserInteractionService) DeleteRegistration(ctx context.Context, eventID int, userID int) error {
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := s.repo.DeleteRegistration(ctx, eventID, userID); err != nil {
		return err
	}

	eventPayload := map[string]string{
		"event_id": strconv.Itoa(eventID),
		"user_id":  strconv.Itoa(userID),
	}

	value, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	if err := s.registrationDeleteW.SendWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 3}, []byte("registration.deleted"), value); err != nil {
		log.Error(ctx, "failed to send registration.deleted", zap.Error(err))
		return err
	}

	log.Info(ctx, "successfully registration.deleted", zap.String("event_id", strconv.Itoa(eventID)), zap.String("user_id", strconv.Itoa(userID)))
	return nil
}

func (s *UserInteractionService) GetRegistrations(ctx context.Context, eventID int, registrations *[]models.ParticipantResp) error {
	return s.repo.GetRegistrations(ctx, eventID, registrations)
}
