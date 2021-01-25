package user

import "github.com/stretchr/testify/mock"

// mockService injects mock dependency into Controller layer.
type mockService struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the UseCase layer.
func (m *mockService) Get(ID string) (*User, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), nil
}

// Create represents the simulated method for the Create feature in the UseCase layer.
func (m *mockService) Create(ID, name, lastname, password string) error {
	args := m.Called(ID, name, lastname, password)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the UseCase layer.
func (m *mockService) Update(u *User) error {
	args := m.Called(u)
	return args.Error(0)
}
