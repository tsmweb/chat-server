package chat

import (
	"encoding/json"
	"errors"
	"github.com/tsmweb/chat-service/config"
	"time"
)

var (
	ErrClosedChannel            = errors.New("closed channel")
)

type ErrorEvent struct {
	HostID    string    `json:"host_id"`
	UserID    string    `json:"user_id,omitempty"`
	Title     string    `json:"title"`
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
}

func NewErrorEvent(userID, title, detail string) *ErrorEvent {
	return &ErrorEvent{
		HostID:    config.HostID(),
		UserID:    userID,
		Title:     title,
		Detail:    detail,
		Timestamp: time.Now().UTC(),
	}
}

func (e ErrorEvent) ToJSON() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return b
}
