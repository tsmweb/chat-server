package group

import "time"

// EventType represents the event type ("delete group", "update group", "add admin"...).
type EventType int

const (
	// EventDeleteGroup represents the group delete event.
	EventDeleteGroup = iota

	// EventUpdateGroup represents the group update event.
	EventUpdateGroup

	// EventAddAdmin represents the admin add event.
	EventAddAdmin

	// EventRemoveAdmin represents the admin remove event.
	EventRemoveAdmin

	// EventAddMember represents the member add event.
	EventAddMember

	// EventRemoveMember represents the member remove event.
	EventRemoveMember
)

var eventTypeText = map[EventType]string{
	EventDeleteGroup:  "DeleteGroup",
	EventUpdateGroup:  "UpdateGroup",
	EventAddAdmin:     "AddAdmin",
	EventRemoveAdmin:  "RemoveAdmin",
	EventAddMember:    "AddMember",
	EventRemoveMember: "RemoveMember",
}

// String return the name of the EventType.
func (e EventType) String() string {
	return eventTypeText[e]
}

// EventTypeText return the name of the EventType.
func EventTypeText(event EventType) string {
	return event.String()
}

// Event represents the events of the domain group.
type Event struct {
	GroupID   string
	MemberID  string
	Event     string
	EventDate time.Time
}

// NewEvent return an instance of Event.
func NewEvent(groupID, memberID string, event EventType) *Event {
	return &Event{
		GroupID:   groupID,
		MemberID:  memberID,
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
