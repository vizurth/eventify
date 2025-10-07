package handler

import (
	"context"
	"eventify/common/logger"
	uipb "eventify/user-interact/api"
	"eventify/user-interact/internal/models"
	"eventify/user-interact/internal/service"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserInteractionHandler struct {
	uipb.UnimplementedUserInteractionServiceServer
	service service.Service
}

func NewUserInteractionHandler(service service.Service) *UserInteractionHandler {
	return &UserInteractionHandler{service: service}
}

// Reviews
func (h *UserInteractionHandler) CreateReview(ctx context.Context, req *uipb.CreateReviewRequest) (*uipb.CreateReviewResponse, error) {
	modelReq := toModelCreateReview(req)

	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := req.ValidateAll(); err != nil {
		log.Error(ctx, "create review:", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.service.CreateNewReviews(ctx, modelReq); err != nil {
		return nil, fmt.Errorf("create review handler: %w", err)
	}
	return &uipb.CreateReviewResponse{Message: "review created"}, nil
}

func (h *UserInteractionHandler) ListReviewsByEvent(ctx context.Context, req *uipb.ListReviewsByEventRequest) (*uipb.ListReviewsByEventResponse, error) {
	var reviews []models.ReviewResp

	if err := h.service.GetCurrentReviewsByEventID(ctx, int(req.GetEventId()), &reviews); err != nil {
		return nil, fmt.Errorf("list reviews handler: %w", err)
	}
	return &uipb.ListReviewsByEventResponse{Reviews: toProtoReviews(reviews)}, nil
}

func (h *UserInteractionHandler) UpdateReview(ctx context.Context, req *uipb.UpdateReviewRequest) (*uipb.UpdateReviewResponse, error) {
	modelReq := toModelUpdateReview(req)

	log := logger.GetOrCreateLoggerFromCtx(ctx)

	if err := req.ValidateAll(); err != nil {
		log.Error(ctx, "create review:", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.service.UpdateReview(ctx, int(req.GetReviewId()), modelReq); err != nil {
		return nil, fmt.Errorf("update review handler: %w", err)
	}
	return &uipb.UpdateReviewResponse{Message: "review updated"}, nil
}

func (h *UserInteractionHandler) DeleteReview(ctx context.Context, req *uipb.DeleteReviewRequest) (*uipb.DeleteReviewResponse, error) {
	if err := h.service.DeleteReview(ctx, int(req.GetReviewId())); err != nil {
		return nil, fmt.Errorf("delete review handler: %w", err)
	}
	return &uipb.DeleteReviewResponse{Message: "review deleted"}, nil
}

// Registrations
func (h *UserInteractionHandler) RegisterForEvent(ctx context.Context, req *uipb.RegisterForEventRequest) (*uipb.RegisterForEventResponse, error) {
	if err := h.service.RegistrationOnEvent(ctx, int(req.GetEventId()), int(req.GetUserId()), req.GetUsername()); err != nil {
		return nil, fmt.Errorf("register on event handler: %w", err)
	}
	return &uipb.RegisterForEventResponse{Message: "registration on"}, nil
}

func (h *UserInteractionHandler) DeleteRegistration(ctx context.Context, req *uipb.DeleteRegistrationRequest) (*uipb.DeleteRegistrationResponse, error) {
	if err := h.service.DeleteRegistration(ctx, int(req.GetEventId()), int(req.GetUserId())); err != nil {
		return nil, fmt.Errorf("delete registration handler: %w", err)
	}
	return &uipb.DeleteRegistrationResponse{Message: "delete registration on"}, nil
}

func (h *UserInteractionHandler) ListRegistrations(ctx context.Context, req *uipb.ListRegistrationsRequest) (*uipb.ListRegistrationsResponse, error) {
	var regs []models.ParticipantResp
	if err := h.service.GetRegistrations(ctx, int(req.GetEventId()), &regs); err != nil {
		return nil, fmt.Errorf("list registrations handler: %w", err)
	}
	return &uipb.ListRegistrationsResponse{Participants: toProtoParticipants(regs)}, nil
}
