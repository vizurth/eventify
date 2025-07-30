package service

import (
	"context"
	"encoding/json"
	"eventify/common/kafka"
	"eventify/event/internal/models"
	"eventify/event/internal/repository"
)

type EventService struct {
	repo     *repository.EventRepository
	producer *kafka.Producer
}

func NewEventService(repo *repository.EventRepository, producer *kafka.Producer) *EventService {
	return &EventService{repo: repo, producer: producer}
}

func (s *EventService) CreateEvent(ctx context.Context, req models.EventReq) error {
	if err := s.repo.CreateEvent(ctx, req); err != nil {
		return err
	}

	// сериализуем данные
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

	// отправляем сообщение
	err = s.producer.SendMessage(ctx, "event.created", value)
	if err != nil {
		return err
	}

	return nil
}

func (s *EventService) GetEvents(ctx context.Context, events *[]models.EventResp) error {
	if err := s.repo.GetEvents(ctx, events); err != nil {
		return err
	}
	return nil
}

func (s *EventService) GetEventByID(ctx context.Context, eventID int, e *models.EventResp) error {
	if err := s.repo.GetEventByID(ctx, eventID, e); err != nil {
		return err
	}
	return nil
}
