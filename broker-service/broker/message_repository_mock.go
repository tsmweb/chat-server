package broker

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/broker-service/broker/message"
)

// mockMessageRepository injects mock message.Repository dependency.
type mockMessageRepository struct {
	mock.Mock
}

// GetAllGroupMembers represents the simulated method for the GetAllGroupMembers
// feature in the message.Repository layer.
func (m *mockMessageRepository) GetAllGroupMembers(ctx context.Context, groupID string) ([]string, error) {
	args := m.Called(ctx, groupID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}

// GetAllMessages represents the simulated method for the GetAllMessages
// feature in the message.Repository layer.
func (m *mockMessageRepository) GetAllMessages(ctx context.Context, userID string) ([]*message.Message, error){
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*message.Message), nil
}

// DeleteAllMessages represents the simulated method for the DeleteAllMessages
// feature in the message.Repository layer.
func (m *mockMessageRepository) DeleteAllMessages(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
