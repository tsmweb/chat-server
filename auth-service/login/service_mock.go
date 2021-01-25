package login

import "github.com/stretchr/testify/mock"

// mockService injects mock dependency into Controller layer.
type mockService struct {
	mock.Mock
}

// Login represents the simulated method for the Login feature in the UseCase layer.
func (m *mockService) Login(ID, password string) (string, error) {
	args := m.Called(ID, password)
	if args.Get(1) != nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), nil
}

// Update represents the simulated method for the Update feature in the UseCase layer.
func (m *mockService) Update(l *Login) error {
	args := m.Called(l)
	return args.Error(0)
}
