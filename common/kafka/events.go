package kafka

import (
	"time"
)

type OrderCreate struct {
	EventType string `json:"event_type"`
	OrderId   int    `json:"order_id"`
	UserId    int    `json:"user_id"`
	Status    string `json:"status"`
}

type OrderStatusUpdate struct {
	EventType string    `json:"event_type"`
	OrderID   int       `json:"order_id"`
	UserId    int       `json:"user_id"`
	CourierId int       `json:"courier_id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderAccept struct {
	EventType string `json:"event_type"`
	OrderID   int    `json:"order_id"`
	Status    string `json:"status"`
}
