package contact

import "github.com/stretchr/testify/mock"

// mockService injects mock dependency into Controller layer.
type mockService struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Service layer.
func (m *mockService) Get(userID, contactID string) (*Contact, error) {
	args := m.Called(userID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Contact), nil
}

// GetAll represents the simulated method for the GetAll feature in the Service layer.
func (m *mockService) GetAll(userID string) ([]*Contact, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Contact), nil
}

// GetPresence represents the simulated method for the GetPresence feature in the Service layer.
func (m *mockService) GetPresence(userID, contactID string) (PresenceType, error) {
	args := m.Called(userID, contactID)
	if args.Error(1) != nil {
		return NotFound, args.Error(1)
	}
	return args.Get(0).(PresenceType), nil
}

// Create represents the simulated method for the Create feature in the Service layer.
func (m *mockService) Create(ID, name, lastname, profileID string) error {
	args := m.Called(ID, name, lastname, profileID)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the Service layer.
func (m *mockService) Update(contact *Contact) error {
	args := m.Called(contact)
	return args.Error(0)
}

// Delete represents the simulated method for the Delete feature in the Service layer.
func (m *mockService) Delete(userID, contactID string) error {
	args := m.Called(userID, contactID)
	return args.Error(0)
}

// Block represents the simulated method for the Block feature in the Service layer.
func (m *mockService) Block(userID, contactID string) error {
	args := m.Called(userID, contactID)
	return args.Error(0)
}

// Unblock represents the simulated method for the Unblock feature in the Service layer.
func (m *mockService) Unblock(userID, contactID string) error {
	args := m.Called(userID, contactID)
	return args.Error(0)
}