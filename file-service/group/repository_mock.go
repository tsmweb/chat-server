package group

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// mockRepository injects mock dependency into UserCase layer.
type mockRepository struct {
	mock.Mock
}

// ExistsGroup represents the simulated method for the ExistsGroup feature in the Repository layer.
func (m *mockRepository) ExistsGroup(ctx context.Context, groupID string) (bool, error) {
	args := m.Called(ctx, groupID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// IsGroupMember represents the simulated method for the IsGroupMember feature in the Repository layer.
func (m *mockRepository) IsGroupMember(ctx context.Context, groupID, userID string) (bool, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// IsGroupAdmin represents the simulated method for the IsGroupAdmin feature in the Repository layer.
func (m *mockRepository) IsGroupAdmin(ctx context.Context, groupID, userID string) (bool, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}
