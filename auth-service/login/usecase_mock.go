package login

import "github.com/stretchr/testify/mock"

// mockLoginUseCase injects mock dependency into Controller layer.
type mockLoginUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Login feature in the UseCase layer.
func (m *mockLoginUseCase) Execute(ID, password string) (string, error) {
	args := m.Called(ID, password)
	if args.Get(1) != nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), nil
}

// mockUpdateUseCase injects mock dependency into Controller layer.
type mockUpdateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Update feature in the UseCase layer.
func (m *mockUpdateUseCase) Execute(l *Login) error {
	args := m.Called(l)
	return args.Error(0)
}
