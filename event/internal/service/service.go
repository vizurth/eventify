package service

import (
	"context"
	"encoding/json"
	"eventify/common/kafka"
	"eventify/event/internal/models"
	"eventify/event/internal/repository"
)

type EventService struct {
	repo      *repository.EventRepository
	producers map[string]*kafka.Producer
}

func NewEventService(repo *repository.EventRepository, producers map[string]*kafka.Producer) *EventService {
	return &EventService{repo: repo, producers: producers}
}

func (s *EventService) CreateEvent(ctx context.Context, req models.EventReq) error {
	if err := s.repo.CreateEvent(ctx, req); err != nil {
		return err
	}

	payload := map[string]string{
		"event_name":     req.Title,
		"event_date":     req.StartTime.Format("2006-01-02 15:04:05"),
		"event_location": req.Location.Address,
		"organizer_name": req.Organizer.Username,
	}
	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if producer := s.producers["event-created"]; producer != nil {
		if err := producer.SendMessage(ctx, "event.created", value); err != nil {
			return err
		}
	}
	return nil
}

func (s *EventService) UpdateEvent(ctx context.Context, eventID int, req models.EventReq) error {
	// Здесь будет логика обновления события
	// Пока просто отправляем в Kafka
	payload := map[string]interface{}{
		"event_id":       eventID,
		"event_name":     req.Title,
		"event_date":     req.StartTime.Format("2006-01-02 15:04:05"),
		"event_location": req.Location.Address,
		"organizer_name": req.Organizer.Username,
	}
	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if producer := s.producers["event-updated"]; producer != nil {
		if err := producer.SendMessage(ctx, "event.updated", value); err != nil {
			return err
		}
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
