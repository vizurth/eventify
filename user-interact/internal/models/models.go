package models

import (
	"database/sql"
	"time"
)

type ReviewReq struct {
	EventID  uint   `json:"event_id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Rating   uint   `json:"rating"`
	Comment  string `json:"comment"`
}

type ReviewResp struct {
	ID        uint         `json:"id"`
	EventID   uint         `json:"event_id"`
	UserID    uint         `json:"user_id"`
	Username  string       `json:"username"`
	Rating    uint         `json:"rating"`
	Comment   string       `json:"comment"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type RegistrationEvent struct {
	EventID  uint   `json:"event_id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type ParticipantResp struct {
	EventID  string `json:"event_id"`
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
