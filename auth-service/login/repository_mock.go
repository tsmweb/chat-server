package login

import "github.com/stretchr/testify/mock"

// mockRepository injects mock dependency into UserCase layer.
type mockRepository struct {
	mock.Mock
}

// Login represents the simulated method for the Login feature in the Repository layer.
func (m *mockRepository) Login(l *Login) (bool, error) {
	args := m.Called(l)
	if args.Get(1) != nil {
		return false, args.Error(1)
	}

	return args.Get(0).(bool), nil
}

// Update represents the simulated method for the Update feature in the
// Repository layer.
func (m *mockRepository) Update(l *Login) (int, error) {
	args := m.Called(l)
	if args.Error(1) != nil {
		return -1, args.Error(1)
	}
	return args.Get(0).(int), nil
}
