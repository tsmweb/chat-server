package user

import "time"

// EventType represents the event type ("block user", "unblock user").
type EventType int

const (
	EventBlockUser   EventType = 0x1
	EventUnblockUser EventType = 0x2
)

// String return the name of the EventType.
func (et EventType) String() (str string) {
	name := func(eventType EventType, name string) bool {
		if et&eventType == 0 {
			return false
		}
		str = name
		return true
	}

	if name(EventBlockUser, "BlockUser") {
		return
	}
	if name(EventUnblockUser, "UnblockUser") {
		return
	}

	return
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

// EventDecoder is a byte slice decoder for user.Event.
type EventDecoder interface {
	Unmarshal(in []byte, evt *Event) error
}

// The EventDecoderFunc type is an adapter to allow the use of ordinary functions as decoders of
// byte slice for user.Event.
// If f is a function with the appropriate signature, EventDecoderFunc(f) is a Decoder that calls f.
type EventDecoderFunc func(in []byte, evt *Event) error

// Unmarshal calls f(in, m).
func (f EventDecoderFunc) Unmarshal(in []byte, evt *Event) error {
	return f(in, evt)
}
