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

// EventDecoder is a byte slice decoder for group.Event.
type EventDecoder interface {
	Unmarshal(in []byte, evt *Event) error
}

// The EventDecoderFunc type is an adapter to allow the use of ordinary functions as decoders of
// byte slice for group.Event.
// If f is a function with the appropriate signature, EventDecoderFunc(f) is a Decoder that calls f.
type EventDecoderFunc func(in []byte, evt *Event) error

// Unmarshal calls f(in, m).
func (f EventDecoderFunc) Unmarshal(in []byte, evt *Event) error {
	return f(in, evt)
}
