package handler

import (
	"context"
	"eventify/common/logger"
	eventpb "eventify/event/api"
	"eventify/event/internal/models"
	"eventify/event/internal/service"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"strconv"
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
	log := logger.GetOrCreateLoggerFromCtx(ctx)
	if err := h.service.CreateEvent(ctx, modelReq); err != nil {
		log.Error(ctx, "create event handler", zap.Error(err))
		return nil, fmt.Errorf("create event handler: %w", err)
	}
	return &eventpb.CreateEventResponse{Message: "event created"}, nil
}

func (h *EventHandler) ListEvents(ctx context.Context, req *eventpb.ListEventsRequest) (*eventpb.ListEventsResponse, error) {
	var events []models.EventResp
	log := logger.GetOrCreateLoggerFromCtx(ctx)
	if err := h.service.GetEvents(ctx, &events); err != nil {
		log.Error(ctx, "list events handler", zap.Error(err))
		return nil, fmt.Errorf("list events handler: %w", err)
	}
	return &eventpb.ListEventsResponse{Events: toProtoEvents(events)}, nil
}

func (h *EventHandler) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.Event, error) {
	var e models.EventResp

	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := h.service.GetEventByID(ctx, int(req.GetId()), &e); err != nil {
		log.Error(ctx, "get events handler", zap.Error(err))
		return nil, fmt.Errorf("get event handler: %w", err)
	}

	if req.UserId != "" {
		userId, err := strconv.Atoi(req.UserId)
		if err != nil {
			log.Error(ctx, "get user id", zap.Error(err))
			return nil, fmt.Errorf("get event: %w", err)
		}
		if err = h.service.CheckUserRegistration(ctx, int(req.GetId()), userId, &e); err != nil {
			log.Error(ctx, "check user", zap.Error(err))
			return nil, fmt.Errorf("get event: %w", err)
		}
	}
	return toProtoEvent(e), nil
}
