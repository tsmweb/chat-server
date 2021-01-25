package user

import (
	"github.com/stretchr/testify/mock"
)

// mockRepository injects mock dependency into UserCase layer.
type mockRepository struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Repository layer.
func (m *mockRepository) Get(ID string) (*User, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), nil
}

// Create represents the simulated method for the Create feature in the
// Repository layer.
func (m *mockRepository) Create(u *User) error {
	args := m.Called(u)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the
// Repository layer.
func (m *mockRepository) Update(u *User) (int, error) {
	args := m.Called(u)
	if args.Error(1) != nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int), nil
}