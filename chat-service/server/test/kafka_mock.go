package test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/go-helper-api/kafka"
)

// mockProducer injects mock kafka.Producer dependency.
type mockProducer struct {
	mock.Mock
}

// Publish represents the simulated method for the Publish feature in the kafka.Producer layer.
func (m *mockProducer) Publish(ctx context.Context, key []byte, value ...[]byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

// Close represents the simulated method for the Close feature in the kafka.Producer layer.
func (m *mockProducer) Close() {}


// mockMessageProducer injects mock kafka.Producer dependency.
type mockMessageProducer struct {
	mock.Mock
	chEvent chan<- kafka.Event
}

// Publish represents the simulated method for the Publish feature in the kafka.Producer layer.
func (m *mockMessageProducer) Publish(ctx context.Context, key []byte, value ...[]byte) error {
	args := m.Called(ctx, key, value)

	for _, v := range value {
		m.chEvent <- kafka.Event{
			Topic: config.KafkaHostTopic(),
			Key: key,
			Value: v,
		}
	}

	return args.Error(0)
}

// Close represents the simulated method for the Close feature in the kafka.Producer layer.
func (m *mockMessageProducer) Close() {}


// mockConsumer injects mock kafka.Consumer dependency.
type mockConsumer struct {
	mock.Mock
}

// Subscribe represents the simulated method for the Subscribe feature in the kafka.Consumer layer.
func (m *mockConsumer) Subscribe(ctx context.Context, callbackFn func(event *kafka.Event, err error)) {
	m.Called(ctx, callbackFn)
}

// Close represents the simulated method for the Close feature in the kafka.Consumer layer.
func (m *mockConsumer) Close() {}
