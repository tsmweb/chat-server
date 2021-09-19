package broker

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

// mockUserRepository injects mock user.Repository dependency.
type mockUserRepository struct {
	mock.Mock
}

// AddUser represents the simulated method for the AddUser feature in the user.Repository layer.
func (m *mockUserRepository) AddUser(ctx context.Context, userID string, serverID string, createAt time.Time) error {
	args := m.Called(ctx, userID, serverID, createAt)
	return args.Error(0)
}

// DeleteUser represents the simulated method for the DeleteUser feature in the user.Repository layer.
func (m *mockUserRepository) DeleteUser(ctx context.Context, userID string, serverID string) error {
	args := m.Called(ctx, userID, serverID)
	return args.Error(0)
}

// GetUserServer represents the simulated method for the GetUserServer feature in the user.Repository layer.
func (m *mockUserRepository) GetUserServer(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

// IsValidUser represents the simulated method for the IsValidUser feature in the user.Repository layer.
func (m *mockUserRepository) IsValidUser(ctx context.Context, fromID string, toID string) (bool, error) {
	args := m.Called(ctx, fromID, toID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// GetAllContactsOnline represents the simulated method for the GetAllContactsOnline
// feature in the user.Repository layer.
func (m *mockUserRepository) GetAllContactsOnline(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}

// GetAllRelationshipsOnline represents the simulated method for the GetAllRelationshipsOnline
// feature in the user.Repository layer.
func (m *mockUserRepository) GetAllRelationshipsOnline(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}
