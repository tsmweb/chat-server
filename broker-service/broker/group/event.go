package group

import "time"

// EventType represents the event type such as EventDeleteGroup, EventUpdateGroup, EventAddAdmin,
// EventRemoveAdmin, EventAddMember and EventRemoveMember.
type EventType int

const (
	EventDeleteGroup  EventType = 0x1
	EventUpdateGroup  EventType = 0x2
	EventAddAdmin     EventType = 0x4
	EventRemoveAdmin  EventType = 0x8
	EventAddMember    EventType = 0x10
	EventRemoveMember EventType = 0x20
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

	if name(EventDeleteGroup, "DeleteGroup") {
		return
	}
	if name(EventUpdateGroup, "UpdateGroup") {
		return
	}
	if name(EventAddAdmin, "AddAdmin") {
		return
	}
	if name(EventRemoveAdmin, "RemoveAdmin") {
		return
	}
	if name(EventAddMember, "AddMember") {
		return
	}
	if name(EventRemoveMember, "RemoveMember") {
		return
	}

	return
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
