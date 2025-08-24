package service

import (
	"context"
	"encoding/json"
	mykafka "eventify/common/kafka"
	"eventify/common/logger"
	"eventify/common/retry"
	"eventify/event/internal/models"
	"eventify/event/internal/repository"
	"go.uber.org/zap"
	"time"
	//"github.com/segmentio/kafka-go"
)

type EventService struct {
	repo          *repository.EventRepository
	eventCreatedW *mykafka.Writer
	//eventUpdatedW *mykafka.Writer
	//eventDeletedW *mykafka.Writer
}

// NewEventService инициализирует сервис, создаёт топики и продюсеров
func NewEventService(ctx context.Context, repo *repository.EventRepository, cfg mykafka.Config) *EventService {
	log := logger.GetOrCreateLoggerFromCtx(ctx)

	// Создаем топики с повторными попытками
	topics := []string{"event-created"}
	for _, topic := range topics {
		if err := mykafka.CreateTopicWithRetry(cfg, topic, 1, 1, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 2}); err != nil {
			log.Error(ctx, "failed to create topic", zap.String("topic", topic), zap.Error(err))
			return nil
		}
		log.Info(ctx, "successfully created topic", zap.String("topic", topic))

	}

	return &EventService{
		repo:          repo,
		eventCreatedW: mykafka.NewWriter(ctx, cfg, "event-created"),
		//eventUpdatedW: mykafka.NewWriter(ctx, cfg, "event-updated"),
		//eventDeletedW: mykafka.NewWriter(ctx, cfg, "event-deleted"),
	}
}

func (s *EventService) CreateEvent(ctx context.Context, req models.EventReq) error {
	log := logger.GetOrCreateLoggerFromCtx(ctx)
	if err := s.repo.CreateEvent(ctx, req); err != nil {
		return err
	}

	eventPayload := map[string]string{
		"event_name":     req.Title,
		"event_date":     req.StartTime.Format("2006-01-02 15:04:05"),
		"event_location": req.Location.Address,
		"organizer_name": req.Organizer.Username,
	}
	value, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	err = s.eventCreatedW.SendWithRetry(ctx, retry.Strategy{Attempts: 3, Delay: time.Second, Backoff: 3}, []byte("event.created"), value)
	if err != nil {
		log.Error(ctx, "failed to write message to kafka", zap.Error(err))
		return err
	}
	log.Info(ctx, "successfully created event", zap.String("event_name", req.Title))
	// отправляем сообщение в топик "event-created"
	return nil
}

func (s *EventService) GetEvents(ctx context.Context, events *[]models.EventResp) error {
	return s.repo.GetEvents(ctx, events)
}

func (s *EventService) GetEventByID(ctx context.Context, eventID int, e *models.EventResp) error {
	return s.repo.GetEventByID(ctx, eventID, e)
}

func (s *EventService) CheckUserRegistration(ctx context.Context, eventID, userID int, e *models.EventResp) error {
	return s.repo.CheckUserRegistration(ctx, eventID, userID, e)
}

func (s *EventService) Close() error {
	if err := s.eventCreatedW.Close(); err != nil {
		return err
	}
	return nil
}
