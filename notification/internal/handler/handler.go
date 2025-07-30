package handler

import (
	"encoding/json"
	"log"
)

type NotificationHandler struct {
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) HandleNotification(eventType string, data []byte) {
	switch eventType {
	case "event.created":
		handleEventCreated(data)
	case "registration.created":
		handleRegistrationCreated(data)
	case "review.created":
		handleReviewCreated(data)
	case "registration.deleted":
		handleRegistrationDeleted(data)
	case "review.updated":
		handleReviewUpdated(data)
	case "review.deleted":
		handleReviewDeleted(data)
	default:
		log.Printf("⚠️ unknown event type: %s", eventType)
	}
}

func handleEventCreated(data []byte) {
	var payload struct {
		EventName     string `json:"event_name"`
		EventDate     string `json:"event_date"`
		EventLocation string `json:"event_location"`
		Organizer     string `json:"organizer_name"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("❌ Invalid event.created payload: %v", err)
		return
	}

	log.Printf("📢 Новое событие: %s\n🗓 %s\n📍 %s\n👤 Организатор: %s",
		payload.EventName, payload.EventDate, payload.EventLocation, payload.Organizer)
}

func handleRegistrationCreated(data []byte) {
	var payload struct {
		UserID    int    `json:"user_id"`
		EventID   int    `json:"event_id"`
		EventName string `json:"event_name"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("❌ Invalid registration.created payload: %v", err)
		return
	}

	log.Printf("✅ Пользователь %d зарегистрировался на событие \"%s\" (ID: %d)",
		payload.UserID, payload.EventName, payload.EventID)
}

func handleRegistrationDeleted(data []byte) {
	var payload struct {
		UserID    int    `json:"user_id"`
		EventID   int    `json:"event_id"`
		EventName string `json:"event_name"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("❌ Invalid registration.deleted payload: %v", err)
		return
	}

	log.Printf("❌ Пользователь %d отменил регистрацию на событие \"%s\" (ID: %d)",
		payload.UserID, payload.EventName, payload.EventID)
}

func handleReviewCreated(data []byte) {
	var payload struct {
		UserID   int    `json:"user_id"`
		EventID  int    `json:"event_id"`
		Rating   int    `json:"rating"`
		Comment  string `json:"comment"`
		UserName string `json:"user_name"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("❌ Invalid review.created payload: %v", err)
		return
	}

	log.Printf("📝 Новый отзыв от %s (ID: %d) на событие %d: %d⭐ — %s",
		payload.UserName, payload.UserID, payload.EventID, payload.Rating, payload.Comment)
}

func handleReviewUpdated(data []byte) {
	var payload struct {
		UserID  int    `json:"user_id"`
		EventID int    `json:"event_id"`
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("❌ Invalid review.updated payload: %v", err)
		return
	}

	log.Printf("🔄 Отзыв пользователя %d на событие %d обновлён: %d⭐ — %s",
		payload.UserID, payload.EventID, payload.Rating, payload.Comment)
}

func handleReviewDeleted(data []byte) {
	var payload struct {
		UserID  int `json:"user_id"`
		EventID int `json:"event_id"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("❌ Invalid review.deleted payload: %v", err)
		return
	}

	log.Printf("🗑 Отзыв пользователя %d на событие %d удалён", payload.UserID, payload.EventID)
}
