package handler

import (
	"context"
	uipb "eventify/user-interact/api"
	"eventify/user-interact/internal/models"
	"eventify/user-interact/internal/service"
)

type UserInteractionHandler struct {
	uipb.UnimplementedUserInteractionServiceServer
	service *service.UserInteractionService
}

func NewUserInteractionHandler(service *service.UserInteractionService) *UserInteractionHandler {
	return &UserInteractionHandler{service: service}
}

// Reviews
func (h *UserInteractionHandler) CreateReview(ctx context.Context, req *uipb.CreateReviewRequest) (*uipb.CreateReviewResponse, error) {
	modelReq := toModelCreateReview(req)
	if err := h.service.CreateNewReviews(ctx, modelReq); err != nil {
		return nil, err
	}
	return &uipb.CreateReviewResponse{Message: "review created"}, nil
}

func (h *UserInteractionHandler) ListReviewsByEvent(ctx context.Context, req *uipb.ListReviewsByEventRequest) (*uipb.ListReviewsByEventResponse, error) {
	var reviews []models.ReviewResp
	if err := h.service.GetCurrentReviewsByEventID(ctx, int(req.GetEventId()), &reviews); err != nil {
		return nil, err
	}
	return &uipb.ListReviewsByEventResponse{Reviews: toProtoReviews(reviews)}, nil
}

func (h *UserInteractionHandler) UpdateReview(ctx context.Context, req *uipb.UpdateReviewRequest) (*uipb.UpdateReviewResponse, error) {
	modelReq := toModelUpdateReview(req)
	if err := h.service.UpdateReview(ctx, int(req.GetReviewId()), modelReq); err != nil {
		return nil, err
	}
	return &uipb.UpdateReviewResponse{Message: "review updated"}, nil
}

func (h *UserInteractionHandler) DeleteReview(ctx context.Context, req *uipb.DeleteReviewRequest) (*uipb.DeleteReviewResponse, error) {
	if err := h.service.DeleteReview(ctx, int(req.GetReviewId())); err != nil {
		return nil, err
	}
	return &uipb.DeleteReviewResponse{Message: "review deleted"}, nil
}

// Registrations
func (h *UserInteractionHandler) RegisterForEvent(ctx context.Context, req *uipb.RegisterForEventRequest) (*uipb.RegisterForEventResponse, error) {
	if err := h.service.RegistrationOnEvent(ctx, int(req.GetEventId()), int(req.GetUserId()), req.GetUsername()); err != nil {
		return nil, err
	}
	return &uipb.RegisterForEventResponse{Message: "registration on"}, nil
}

func (h *UserInteractionHandler) DeleteRegistration(ctx context.Context, req *uipb.DeleteRegistrationRequest) (*uipb.DeleteRegistrationResponse, error) {
	if err := h.service.DeleteRegistration(ctx, int(req.GetEventId()), int(req.GetUserId())); err != nil {
		return nil, err
	}
	return &uipb.DeleteRegistrationResponse{Message: "delete registration on"}, nil
}

func (h *UserInteractionHandler) ListRegistrations(ctx context.Context, req *uipb.ListRegistrationsRequest) (*uipb.ListRegistrationsResponse, error) {
	var regs []models.ParticipantResp
	if err := h.service.GetRegistrations(ctx, int(req.GetEventId()), &regs); err != nil {
		return nil, err
	}
	return &uipb.ListRegistrationsResponse{Participants: toProtoParticipants(regs)}, nil
}
