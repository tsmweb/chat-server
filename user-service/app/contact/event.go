package contact

import "time"

// EventType represents the event type ("block user", "unblock user").
type EventType int

const (
	// EventBlockUser represents the block user event.
	EventBlockUser = iota

	// EventUnblockUser represents the unblock user event.
	EventUnblockUser
)

var eventTypeText = map[EventType]string{
	EventBlockUser:   "BlockUser",
	EventUnblockUser: "UnblockUser",
}

// String return the name of the EventType.
func (e EventType) String() string {
	return eventTypeText[e]
}

// EventTypeText return the name of the EventType.
func EventTypeText(event EventType) string {
	return event.String()
}

// Event represents the events of the domain contact.
type Event struct {
	UserID    string
	ContactID string
	Event     string
	EventDate time.Time
}

// NewEvent return an instance of Event.
func NewEvent(userID, contactID string, event EventType) *Event {
	return &Event{
		UserID:    userID,
		ContactID: contactID,
		Event:     event.String(),
		EventDate: time.Now().UTC(),
	}
}

// EventEncoder is a Event encoder for byte slice
type EventEncoder interface {
	Marshal(e *Event) ([]byte, error)
}

// The EventEncoderFunc type is an adapter to allow the use of ordinary functions as encoders of Event for byte slice.
// If f is a function with the appropriate signature, EventEncoderFunc(f) is a Encoder that calls f.
type EventEncoderFunc func(e *Event) ([]byte, error)

// Marshal calls f(e).
func (f EventEncoderFunc) Marshal(e *Event) ([]byte, error) {
	return f(e)
}
