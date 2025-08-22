package service

import (
	"context"
	"encoding/json"
	mykafka "eventify/common/kafka"
	"eventify/common/logger"
	"eventify/event/internal/models"
	"eventify/event/internal/repository"

	"github.com/segmentio/kafka-go"
)

type EventService struct {
	repo          *repository.EventRepository
	eventCreatedW *kafka.Writer
	eventUpdatedW *kafka.Writer
	eventDeletedW *kafka.Writer
}

// NewEventService инициализирует сервис, создаёт топики и продюсеров
func NewEventService(ctx context.Context, repo *repository.EventRepository, cfg mykafka.Config) *EventService {
	log := logger.GetOrCreateLoggerFromCtx(ctx)
	if err := mykafka.CreateTopicWithRetry(cfg, "event-created", 1, 1); err != nil {
		log.Error(ctx, "failed to create topic event-created")
		return nil
	}
	if err := mykafka.CreateTopicWithRetry(cfg, "event-updated", 1, 1); err != nil {
		log.Error(ctx, "failed to create topic event-updated")
		return nil
	}
	if err := mykafka.CreateTopicWithRetry(cfg, "event-deleted", 1, 1); err != nil {
		log.Error(ctx, "failed to create topic event-deleted")
		return nil
	}

	return &EventService{
		repo:          repo,
		eventCreatedW: mykafka.NewWriter(ctx, cfg, "event-created"),
		eventUpdatedW: mykafka.NewWriter(ctx, cfg, "event-updated"),
		eventDeletedW: mykafka.NewWriter(ctx, cfg, "event-deleted"),
	}
}

func (s *EventService) CreateEvent(ctx context.Context, req models.EventReq) error {
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

	// отправляем сообщение в топик "event-created"
	return s.eventCreatedW.WriteMessages(ctx, kafka.Message{
		Key:   []byte("event.created"),
		Value: value,
	})
}

func (s *EventService) GetEvents(ctx context.Context, events *[]models.EventResp) error {
	return s.repo.GetEvents(ctx, events)
}

func (s *EventService) GetEventByID(ctx context.Context, eventID int, e *models.EventResp) error {
	return s.repo.GetEventByID(ctx, eventID, e)
}

func (s *EventService) Close() error {
	if err := s.eventCreatedW.Close(); err != nil {
		return err
	}
	if err := s.eventUpdatedW.Close(); err != nil {
		return err
	}
	if err := s.eventDeletedW.Close(); err != nil {
		return err
	}
	return nil
}
