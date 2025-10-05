package repository

import (
	"context"
	"eventify/event/internal/models"
)

type Repository interface {
	CreateEvent(ctx context.Context, req models.EventReq) error
	GetEvents(ctx context.Context, events *[]models.EventResp) error
	GetEventByID(ctx context.Context, eventID int, e *models.EventResp) error
	CheckUserRegistration(ctx context.Context, eventID, userID int, e *models.EventResp) error
}
