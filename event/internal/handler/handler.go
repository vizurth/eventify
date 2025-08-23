package handler

import (
	"context"
	eventpb "eventify/event/api"
	"eventify/event/internal/models"
	"eventify/event/internal/service"
)

type EventHandler struct {
	eventpb.UnimplementedEventServiceServer
	service *service.EventService
}

func NewEventHandler(s *service.EventService) *EventHandler {
	return &EventHandler{service: s}
}

func (h *EventHandler) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	modelReq := toModelCreate(req)
	if err := h.service.CreateEvent(ctx, modelReq); err != nil {
		return nil, err
	}
	return &eventpb.CreateEventResponse{Message: "event created"}, nil
}

func (h *EventHandler) ListEvents(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) {
	var events []models.EventResp
	if err := h.service.GetEvents(ctx, &events); err != nil {
		return nil, err
	}
	return &eventpb.ListEventsResponse{Events: toProtoEvents(events)}, nil
}

func (h *EventHandler) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.Event, error) {
	var e models.EventResp
	if err := h.service.GetEventByID(ctx, int(req.GetId()), &e); err != nil {
		return nil, err
	}
	return toProtoEvent(e), nil
}
