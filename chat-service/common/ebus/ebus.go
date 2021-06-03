package ebus

import (
	"context"
	"time"
)

// Event structure that represents a event streaming.
type Event struct {
	Key   string
	Value []byte
	Time  time.Time
}

// EBus it is an abstraction to send and consume events from an
// event streaming service.
type EBus interface {

	// Dispatch produces and sends an event for a topic.
	// The context passed as first argument may also be used to asynchronously
	// cancel the operation.
	Dispatch(ctx context.Context, topic string, key string, value []byte) error

	// Subscribe consumes the events of a topic and passes the event to the
	// informed callback function. The method call blocks until an error occurs.
	// The program may also specify a context to asynchronously cancel the blocking operation.
	Subscribe(ctx context.Context, groupID, topic string, callbackFn func(event *Event, err error))
}
