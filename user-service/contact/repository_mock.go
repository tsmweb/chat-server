package contact

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

// mockRepository injects mock dependency into UserCase layer.
type mockRepository struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Repository layer.
func (m *mockRepository) Get(ctx context.Context, profileID, contactID string) (*Contact, error) {
	args := m.Called(ctx, profileID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Contact), nil
}

// GetAll represents the simulated method for the GetAll feature in the Repository layer.
func (m *mockRepository) GetAll(ctx context.Context, profileID string) ([]*Contact, error) {
	args := m.Called(ctx, profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Contact), nil
}

// ExistsUser represents the simulated method for the ExistsUser feature in the Repository layer.
func (m *mockRepository) ExistsUser(ctx context.Context, ID string) (bool, error) {
	args := m.Called(ctx, ID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// GetPresence represents the simulated method for the GetPresence feature in the Repository layer.
func (m *mockRepository) GetPresence(ctx context.Context, profileID, contactID string) (PresenceType, error) {
	args := m.Called(ctx, profileID, contactID)
	if args.Error(1) != nil {
		return NotFound, args.Error(1)
	}
	return args.Get(0).(PresenceType), nil
}

// Create represents the simulated method for the Create feature in the Repository layer.
func (m *mockRepository) Create(ctx context.Context, c *Contact) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the Repository layer.
func (m *mockRepository) Update(ctx context.Context, c *Contact) (bool, error) {
	args := m.Called(ctx, c)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// Delete represents the simulated method for the Delete feature in the Repository layer.
func (m *mockRepository) Delete(ctx context.Context, userID, contactID string) (bool, error) {
	args := m.Called(ctx, userID, contactID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// Block represents the simulated method for the Block feature in the Repository layer.
func (m *mockRepository) Block(ctx context.Context, profileID, blockedUserID string, createdAt time.Time) error {
	args := m.Called(ctx, profileID, blockedUserID, createdAt)
	return args.Error(0)
}

// Unblock represents the simulated method for the Unblock feature in the Repository layer.
func (m *mockRepository) Unblock(ctx context.Context, profileID, blockedUserID string) (bool, error) {
	args := m.Called(ctx, profileID, blockedUserID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}
