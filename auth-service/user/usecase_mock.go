package user

import "github.com/stretchr/testify/mock"

// mockGetUseCase injects mock dependency into Controller layer.
type mockGetUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Get feature in the UseCase layer.
func (m *mockGetUseCase) Execute(ID string) (*User, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), nil
}

// mockCreateUseCase injects mock dependency into Controller layer.
type mockCreateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Create feature in the UseCase layer.
func (m *mockCreateUseCase) Execute(ID, name, lastname, password string) error {
	args := m.Called(ID, name, lastname, password)
	return args.Error(0)
}

// mockUpdateUseCase injects mock dependency into Controller layer.
type mockUpdateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Update feature in the UseCase layer.
func (m *mockUpdateUseCase) Execute(u *User) error {
	args := m.Called(u)
	return args.Error(0)
}