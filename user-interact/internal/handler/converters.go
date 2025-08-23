package handler

import (
	uipb "eventify/user-interact/api"
	"eventify/user-interact/internal/models"
	"time"
)

func toProtoReviews(src []models.ReviewResp) []*uipb.Review {
	out := make([]*uipb.Review, 0, len(src))
	for _, r := range src {
		out = append(out, toProtoReview(r))
	}
	return out
}

func toProtoReview(r models.ReviewResp) *uipb.Review {
	var updated string
	if r.UpdatedAt.Valid {
		updated = r.UpdatedAt.Time.Format(time.RFC3339)
	}
	return &uipb.Review{
		Id:        uint64(r.ID),
		EventId:   uint64(r.EventID),
		UserId:    uint64(r.UserID),
		Username:  r.Username,
		Rating:    uint32(r.Rating),
		Comment:   r.Comment,
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updated,
	}
}

func toProtoParticipants(src []models.ParticipantResp) []*uipb.Participant {
	out := make([]*uipb.Participant, 0, len(src))
	for _, p := range src {
		out = append(out, &uipb.Participant{Id: uint64(p.ID), Username: p.Username})
	}
	return out
}

// toModelCreateReview converts CreateReviewRequest into internal ReviewReq
func toModelCreateReview(req *uipb.CreateReviewRequest) models.ReviewReq {
	return models.ReviewReq{
		EventID:  uint(req.GetEventId()),
		UserID:   uint(req.GetUserId()),
		Username: req.GetUsername(),
		Rating:   uint(req.GetRating()),
		Comment:  req.GetComment(),
	}
}

// toModelUpdateReview converts UpdateReviewRequest into internal ReviewReq
func toModelUpdateReview(req *uipb.UpdateReviewRequest) models.ReviewReq {
	return models.ReviewReq{
		EventID:  uint(req.GetEventId()),
		UserID:   uint(req.GetUserId()),
		Username: req.GetUsername(),
		Rating:   uint(req.GetRating()),
		Comment:  req.GetComment(),
	}
}
