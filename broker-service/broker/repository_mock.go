package broker

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/broker-service/broker/message"
	"time"
)

// mockUserRepository injects mock user.Repository dependency.
type mockUserRepository struct {
	mock.Mock
}

// AddUserPresence represents the simulated method for the AddUserPresence feature in the
// user.Repository layer.
func (m *mockUserRepository) AddUserPresence(ctx context.Context, userID string, serverID string,
	createAt time.Time) error {
	args := m.Called(ctx, userID, serverID, createAt)
	return args.Error(0)
}

// RemoveUserPresence represents the simulated method for the RemoveUserPresence feature in the
// user.Repository layer.
func (m *mockUserRepository) RemoveUserPresence(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// UpdateUserPresenceCache represents the simulated method for the UpdateUserPresenceCache
// feature in the user.Repository layer.
func (m *mockUserRepository) UpdateUserPresenceCache(ctx context.Context, userID string,
	serverID string, status string) error {
	args := m.Called(ctx, userID, serverID, status)
	return args.Error(0)
}

// GetUserServer represents the simulated method for the GetUserServer feature in the
// user.Repository layer.
func (m *mockUserRepository) GetUserServer(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

// IsValidUser represents the simulated method for the IsValidUser feature in the
// user.Repository layer.
func (m *mockUserRepository) IsValidUser(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// IsBlockedUser represents the simulated method for the IsBlockedUser feature in the
//user.Repository layer.
func (m *mockUserRepository) IsBlockedUser(ctx context.Context, fromID string,
	toID string) (bool, error) {
	args := m.Called(ctx, fromID, toID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// UpdateBlockedUserCache represents the simulated method for the UpdateBlockedUserCache feature
// in the user.Repository layer.
func (m *mockUserRepository) UpdateBlockedUserCache(ctx context.Context, userID string,
	blockedUserID string, blocked bool) error {
	args := m.Called(ctx, userID, blockedUserID, blocked)
	return args.Error(0)
}

// GetAllContactsOnline represents the simulated method for the GetAllContactsOnline
// feature in the user.Repository layer.
func (m *mockUserRepository) GetAllContactsOnline(ctx context.Context,
	userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}

// GetAllRelationshipsOnline represents the simulated method for the GetAllRelationshipsOnline
// feature in the user.Repository layer.
func (m *mockUserRepository) GetAllRelationshipsOnline(ctx context.Context,
	userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}

// mockMessageRepository injects mock message.Repository dependency.
type mockMessageRepository struct {
	mock.Mock
}

// GetAllGroupMembers represents the simulated method for the GetAllGroupMembers
// feature in the message.Repository layer.
func (m *mockMessageRepository) GetAllGroupMembers(ctx context.Context,
	groupID string) ([]string, error) {
	args := m.Called(ctx, groupID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}

// RemoveGroupFromCache represents the simulated method for the RemoveGroupFromCache
// feature in the message.Repository layer.
func (m *mockMessageRepository) RemoveGroupFromCache(ctx context.Context, groupID string) error {
	args := m.Called(ctx, groupID)
	return args.Error(0)
}

// AddGroupMemberToCache represents the simulated method for the AddGroupMemberToCache
// feature in the message.Repository layer.
func (m *mockMessageRepository) AddGroupMemberToCache(ctx context.Context, groupID,
	memberID string) error {
	args := m.Called(ctx, groupID, memberID)
	return args.Error(0)
}

// RemoveGroupMemberFromCache represents the simulated method for the RemoveGroupMemberFromCache
// feature in the message.Repository layer.
func (m *mockMessageRepository) RemoveGroupMemberFromCache(ctx context.Context, groupID,
	memberID string) error {
	args := m.Called(ctx, groupID, memberID)
	return args.Error(0)
}

// GetAllMessages represents the simulated method for the GetAllMessages
// feature in the message.Repository layer.
func (m *mockMessageRepository) GetAllMessages(ctx context.Context,
	userID string) ([]*message.Message, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*message.Message), nil
}

// AddMessage represents the simulated method for the AddMessage
// feature in the message.Repository layer.
func (m *mockMessageRepository) AddMessage(ctx context.Context, msg message.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

// DeleteAllMessages represents the simulated method for the DeleteAllMessages
// feature in the message.Repository layer.
func (m *mockMessageRepository) DeleteAllMessages(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
