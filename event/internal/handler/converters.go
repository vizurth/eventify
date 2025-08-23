package handler

import (
	eventpb "eventify/event/api"
	"eventify/event/internal/models"
	"time"
)

func toModelCreate(req *eventpb.CreateEventRequest) models.EventReq {
	start, _ := time.Parse(time.RFC3339, req.GetStartTime())
	end, _ := time.Parse(time.RFC3339, req.GetEndTime())
	participants := make([]models.Participant, 0, len(req.GetParticipants()))
	for _, p := range req.GetParticipants() {
		participants = append(participants, models.Participant{ID: uint(p.GetId()), Username: p.GetUsername()})
	}
	return models.EventReq{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Category:    req.GetCategory(),
		Location: models.Location{
			City:    req.GetLocation().GetCity(),
			Venue:   req.GetLocation().GetVenue(),
			Address: req.GetLocation().GetAddress(),
		},
		StartTime: start,
		EndTime:   end,
		Organizer: models.Organizer{
			ID:       uint(req.GetOrganizer().GetId()),
			Username: req.GetOrganizer().GetUsername(),
			Email:    req.GetOrganizer().GetEmail(),
		},
		Participants: participants,
		Status:       req.GetStatus(),
	}
}

func toProtoEvents(src []models.EventResp) []*eventpb.Event {
	out := make([]*eventpb.Event, 0, len(src))
	for _, e := range src {
		out = append(out, toProtoEvent(e))
	}
	return out
}

func toProtoEvent(e models.EventResp) *eventpb.Event {
	parts := make([]*eventpb.Participant, 0, len(e.Participants))
	for _, p := range e.Participants {
		parts = append(parts, &eventpb.Participant{Id: uint64(p.ID), Username: p.Username})
	}
	return &eventpb.Event{
		Id:          uint64(e.ID),
		Title:       e.Title,
		Description: e.Description,
		Category:    e.Category,
		Location:    &eventpb.Location{City: e.Location.City, Venue: e.Location.Venue, Address: e.Location.Address},
		StartTime:   e.StartTime.Format(time.RFC3339),
		EndTime:     e.EndTime.Format(time.RFC3339),
		Organizer:   &eventpb.Organizer{Id: uint64(e.Organizer.ID), Username: e.Organizer.Username, Email: e.Organizer.Email},
		Participants: parts,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt.Format(time.RFC3339),
	}
} 