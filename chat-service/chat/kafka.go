package chat

import (
	"context"
	"time"
)

// KafkaEvent structure that represents a Kafka event.
type KafkaEvent struct {
	Topic  string
	Key    []byte
	Value  []byte
	Time   time.Time
}

// Kafka it is an abstraction to send and consume events from an
// event kafka service.
type Kafka interface {

	// NewProducer creates a new KafkaProducer to produce events on a topic.
	NewProducer(topic string) KafkaProducer

	// NewConsumer creates a new KafkaConsumer to consume events from a topic.
	NewConsumer(groupID, topic string) KafkaConsumer

	// Debug enables logging of incoming events.
	Debug(debug bool)
}

// KafkaProducer provide methods for producing events for a given topic.
type KafkaProducer interface {
	// Publish produces and sends an event for a kafka topic.
	// The context passed as first argument may also be used to asynchronously
	// cancel the operation.
	Publish(ctx context.Context, key, value []byte) error

	// Close flushes pending writes, and waits for all writes to complete before
	// returning.
	Close()
}

// KafkaConsumer provide methods for consuming events on a given topic.
type KafkaConsumer interface {
	// Subscribe consumes the events of a topic and passes the event to the
	// informed callback function. The method call blocks until an error occurs.
	// The program may also specify a context to asynchronously cancel the blocking operation.
	Subscribe(ctx context.Context, callback func(event *KafkaEvent, err error))

	// Close closes the stream, preventing the program from reading any more
	// events from it.
	Close()
}
