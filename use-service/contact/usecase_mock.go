package contact

import "github.com/stretchr/testify/mock"

// mockGetUseCase injects mock dependency into Controller layer.
type mockGetUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Get feature in the UseCase layer.
func (m *mockGetUseCase) Execute(profileID, contactID string) (*Contact, error) {
	args := m.Called(profileID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Contact), nil
}

// mockGetAllUseCase injects mock dependency into Controller layer.
type mockGetAllUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the GetAll feature in the UseCase layer.
func (m *mockGetAllUseCase) Execute(profileID string) ([]*Contact, error) {
	args := m.Called(profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Contact), nil
}

// mockGetPresenceUseCase injects mock dependency into Controller layer.
type mockGetPresenceUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the GetPresence feature in the UseCase layer.
func (m *mockGetPresenceUseCase) Execute(profileID, contactID string) (PresenceType, error) {
	args := m.Called(profileID, contactID)
	if args.Error(1) != nil {
		return NotFound, args.Error(1)
	}
	return args.Get(0).(PresenceType), nil
}

// mockCreateUseCase injects mock dependency into Controller layer.
type mockCreateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Create feature in the UseCase layer.
func (m *mockCreateUseCase) Execute(ID, name, lastname, profileID string) error {
	args := m.Called(ID, name, lastname, profileID)
	return args.Error(0)
}

// mockUpdateUseCase injects mock dependency into Controller layer.
type mockUpdateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Update feature in the UseCase layer.
func (m *mockUpdateUseCase) Execute(contact *Contact) error {
	args := m.Called(contact)
	return args.Error(0)
}

// mockDeleteUseCase injects mock dependency into Controller layer.
type mockDeleteUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Delete feature in the UseCase layer.
func (m *mockDeleteUseCase) Execute(contact *Contact) error {
	args := m.Called(contact)
	return args.Error(0)
}

// mockBlockUseCase injects mock dependency into Controller layer.
type mockBlockUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Block feature in the UseCase layer.
func (m *mockBlockUseCase) Execute(profileID, contactID string) error {
	args := m.Called(profileID, contactID)
	return args.Error(0)
}

// mockUnblockUseCase injects mock dependency into Controller layer.
type mockUnblockUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Unblock feature in the UseCase layer.
func (m *mockUnblockUseCase) Execute(profileID, contactID string) error {
	args := m.Called(profileID, contactID)
	return args.Error(0)
}
