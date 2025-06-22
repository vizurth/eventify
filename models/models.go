package models

import (
	"database/sql"
	"time"
)

type LoginStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Location struct {
	City    string `json:"city"`
	Venue   string `json:"venue"`
	Address string `json:"address"`
}

type Organizer struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Participant struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type ParticipantResp struct {
	EventID string `json:"event_id"`

	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type EventReq struct {
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Category     string        `json:"category"`
	Location     Location      `json:"location"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Organizer    Organizer     `json:"organizer"`
	Participants []Participant `json:"participants"`
	Status       string        `json:"status"`
}
type EventResp struct {
	ID           uint          `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Category     string        `json:"category"`
	Location     Location      `json:"location"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Organizer    Organizer     `json:"organizer"`
	Participants []Participant `json:"participants"`
	Status       string        `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
}

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
