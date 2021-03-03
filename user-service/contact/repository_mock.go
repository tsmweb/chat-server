package contact

import (
	"github.com/stretchr/testify/mock"
	"time"
)

// mockRepository injects mock dependency into UserCase layer.
type mockRepository struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Repository layer.
func (m *mockRepository) Get(profileID, contactID string) (*Contact, error) {
	args := m.Called(profileID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Contact), nil
}

// GetAll represents the simulated method for the GetAll feature in the Repository layer.
func (m *mockRepository) GetAll(profileID string) ([]*Contact, error) {
	args := m.Called(profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Contact), nil
}

// ExistsUser represents the simulated method for the ExistsUser feature in the Repository layer.
func (m *mockRepository) ExistsUser(ID string) (bool, error) {
	args := m.Called(ID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// GetPresence represents the simulated method for the GetPresence feature in the Repository layer.
func (m *mockRepository) GetPresence(profileID, contactID string) (PresenceType, error) {
	args := m.Called(profileID, contactID)
	if args.Error(1) != nil {
		return NotFound, args.Error(1)
	}
	return args.Get(0).(PresenceType), nil
}

// Create represents the simulated method for the Create feature in the Repository layer.
func (m *mockRepository) Create(c *Contact) error {
	args := m.Called(c)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the Repository layer.
func (m *mockRepository) Update(c *Contact) (int, error) {
	args := m.Called(c)
	if args.Error(1) != nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), nil
}

// Delete represents the simulated method for the Delete feature in the Repository layer.
func (m *mockRepository) Delete(userID, contactID string) (int, error) {
	args := m.Called(userID, contactID)
	if args.Error(1) != nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), nil
}

// Block represents the simulated method for the Block feature in the Repository layer.
func (m *mockRepository) Block(profileID, blockedUserID string, createdAt time.Time) error {
	args := m.Called(profileID, blockedUserID, createdAt)
	return args.Error(0)
}

// Unblock represents the simulated method for the Unblock feature in the Repository layer.
func (m *mockRepository) Unblock(profileID, blockedUserID string) (bool, error) {
	args := m.Called(profileID, blockedUserID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}