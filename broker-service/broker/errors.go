package broker

import (
	"encoding/json"
	"errors"
	"github.com/tsmweb/broker-service/config"
	"strings"
	"time"
)

var (
	ErrClosedChannel = errors.New("closed channel")
)

type ErrorEvent struct {
	HostID    string    `json:"host_id"`
	UserID    string    `json:"user_id,omitempty"`
	Title     string    `json:"title"`
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
}

func NewErrorEvent(userID, title, detail string) *ErrorEvent {
	if strings.TrimSpace(userID) == "" {
		userID = config.HostID()
	}

	return &ErrorEvent{
		HostID:    config.HostID(),
		UserID:    userID,
		Title:     title,
		Detail:    detail,
		Timestamp: time.Now().UTC(),
	}
}

func (e *ErrorEvent) ToJSON() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return b
}

// ErrorEventEncoder is a ErrorEvent encoder for byte slice.
type ErrorEventEncoder interface {
	Marshal(e *ErrorEvent) ([]byte, error)
}

// The ErrorEventEncoderFunc type is an adapter to allow the use of ordinary functions as encoders of ErrorEvent for byte slice.
// If f is a function with the appropriate signature, EncoderFunc(f) is a Encoder that calls f.
type ErrorEventEncoderFunc func(e *ErrorEvent) ([]byte, error)

// Marshal calls f(m).
func (f ErrorEventEncoderFunc) Marshal(e *ErrorEvent) ([]byte, error) {
	return f(e)
}