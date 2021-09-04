package group

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// mockProducer injects mock kafka.Producer dependency.
type mockProducer struct {
	mock.Mock
}

// Publish represents the simulated method for the Publish feature in the kafka.Producer layer.
func (m *mockProducer) Publish(ctx context.Context, key, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

// Close represents the simulated method for the Close feature in the kafka.Producer layer.
func (m *mockProducer) Close() {}
