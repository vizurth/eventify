package service

import (
	"context"
	"vizurth/eventify/event/internal/repository"
	"vizurth/eventify/models"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, req models.EventReq) error {
	if err := s.repo.CreateEvent(ctx, req); err != nil {
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
