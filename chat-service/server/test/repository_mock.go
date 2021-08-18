package test

import (
	"github.com/stretchr/testify/mock"
	"time"
)

// mockProducer injects mock Repository dependency.
type mockRepository struct {
	mock.Mock
}

// AddUserOnline represents the simulated method for the AddUserOnline feature
// in the Repository layer.
func (m *mockRepository) AddUserOnline(userID string, host string, createAt time.Time) error {
	args := m.Called(userID, host, createAt)
	return args.Error(0)
}

// DeleteUserOnline represents the simulated method for the DeleteUserOnline feature
// in the Repository layer.
func (m *mockRepository) DeleteUserOnline(userID string) error {
	args := m.Called(userID)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

// IsValidUser represents the simulated method for the IsValidUser feature in the Repository layer.
func (m *mockRepository) IsValidUser(fromID string, toID string) (bool, error) {
	args := m.Called(fromID, toID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// GetGroupMembers represents the simulated method for the GetGroupMembers feature
// in the Repository layer.
func (m *mockRepository) GetGroupMembers(groupID string) ([]string, error) {
	args := m.Called(groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}

// GetUserContactsOnline represents the simulated method for the GetUserContactsOnline feature
// in the Repository layer.
func (m *mockRepository) GetUserContactsOnline(userID string) ([]string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}

// GetContactsWithUserOnline represents the simulated method for the GetContactsWithUserOnline feature
// in the Repository layer.
func (m *mockRepository) GetContactsWithUserOnline(userID string) ([]string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}
