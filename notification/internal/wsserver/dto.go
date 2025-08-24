package wsserver

import (
	"github.com/gorilla/websocket"
	"time"
)

type wsMessage struct {
	IpAddress string `json:"ip_address"`
	Message   string `json:"message"`
	Time      string `json:"time"`
}

type NotificationMessage struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	EventDate string `json:"event_date,omitempty"`
	Location  string `json:"location,omitempty"`
	Organizer string `json:"organizer,omitempty"`
	Timestamp string `json:"timestamp"`
}

type ClientConnection struct {
	ID       string
	Send     chan []byte
	Conn     *websocket.Conn
	LastPing time.Time
}
