package broker

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/kafka"
)

// mockKafka injects mock kafka.Kafka dependency.
type mockKafka struct {
	mock.Mock
}

// NewProducer represents the simulated method for the NewProducer feature in the kafka.Kafka layer.
func (m *mockKafka) NewProducer(topic string) kafka.Producer {
	args := m.Called(topic)
	return args.Get(0).(kafka.Producer)
}

// NewConsumer represents the simulated method for the NewConsumer feature in the kafka.Kafka layer.
func (m *mockKafka) NewConsumer(groupID, topic string) kafka.Consumer {
	args := m.Called(groupID, topic)
	return args.Get(0).(kafka.Consumer)
}

// Debug represents the simulated method for the Debug feature in the kafka.Kafka layer.
func (m *mockKafka) Debug(debug bool) {}


// mockProducer injects mock kafka.Producer dependency.
type mockProducer struct {
	mock.Mock
}

// Publish represents the simulated method for the Publish feature in the kafka.Producer layer.
func (m *mockProducer) Publish(ctx context.Context, key []byte, values ...[]byte) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

// Close represents the simulated method for the Close feature in the kafka.Producer layer.
func (m *mockProducer) Close() {}


// mockConsumer injects mock kafka.Consumer dependency.
type mockConsumer struct {
	mock.Mock
}

// Subscribe represents the simulated method for the Subscribe feature in the kafka.Consumer layer.
func (m *mockConsumer) Subscribe(ctx context.Context, callbackFn func(event *kafka.Event, err error)) {}

// Close represents the simulated method for the Close feature in the kafka.Consumer layer.
func (m *mockConsumer) Close() {}
