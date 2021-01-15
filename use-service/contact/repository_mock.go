package contact

import "github.com/stretchr/testify/mock"

// mockRepository injects mock dependency into UserCase layer.
type mockRepository struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Repository layer.
func (m *mockRepository) Get(profileID, contactID string) (Contact, error) {
	c := Contact{}
	args := m.Called(profileID, contactID)
	if args.Get(0) == nil {
		return c, args.Error(1)
	}
	return args.Get(0).(Contact), nil
}

// GetAll represents the simulated method for the GetAll feature in the Repository layer.
func (m *mockRepository) GetAll(profileID string) ([]Contact, error) {
	var c []Contact
	args := m.Called(profileID)
	if args.Get(0) == nil {
		return c, args.Error(1)
	}
	return args.Get(0).([]Contact), nil
}

// ExistsProfile represents the simulated method for the ExistsProfile feature in the Repository layer.
func (m *mockRepository) ExistsProfile(ID string) (bool, error) {
	args := m.Called(ID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// Create represents the simulated method for the Create feature in the Repository layer.
func (m *mockRepository) Create(c Contact) error {
	args := m.Called(c)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the Repository layer.
func (m *mockRepository) Update(c Contact) error {
	args := m.Called(c)
	return args.Error(0)
}

// Delete represents the simulated method for the Delete feature in the Repository layer.
func (m *mockRepository) Delete(c Contact) error {
	args := m.Called(c)
	return args.Error(0)
}