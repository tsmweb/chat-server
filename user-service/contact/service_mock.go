package contact

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// mockService injects mock dependency into Controller layer.
type mockService struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Service layer.
func (m *mockService) Get(ctx context.Context, userID, contactID string) (*Contact, error) {
	args := m.Called(ctx, userID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Contact), nil
}

// GetAll represents the simulated method for the GetAll feature in the Service layer.
func (m *mockService) GetAll(ctx context.Context, userID string) ([]*Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Contact), nil
}

// GetPresence represents the simulated method for the GetPresence feature in the Service layer.
func (m *mockService) GetPresence(ctx context.Context, userID, contactID string) (PresenceType, error) {
	args := m.Called(ctx, userID, contactID)
	if args.Error(1) != nil {
		return NotFound, args.Error(1)
	}
	return args.Get(0).(PresenceType), nil
}

// Create represents the simulated method for the Create feature in the Service layer.
func (m *mockService) Create(ctx context.Context, ID, name, lastname, profileID string) error {
	args := m.Called(ctx, ID, name, lastname, profileID)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the Service layer.
func (m *mockService) Update(ctx context.Context, contact *Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

// Delete represents the simulated method for the Delete feature in the Service layer.
func (m *mockService) Delete(ctx context.Context, userID, contactID string) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

// Block represents the simulated method for the Block feature in the Service layer.
func (m *mockService) Block(ctx context.Context, userID, contactID string) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

// Unblock represents the simulated method for the Unblock feature in the Service layer.
func (m *mockService) Unblock(ctx context.Context, userID, contactID string) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}
