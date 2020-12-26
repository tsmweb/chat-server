package profile

import "github.com/stretchr/testify/mock"

// mockGetUseCase injects mock dependency into Controller layer.
type mockGetUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Get feature in the UseCase layer.
func (m *mockGetUseCase) Execute(ID string) (Profile, error) {
	p := Profile{}
	args := m.Called(ID)
	if args.Get(0) == nil {
		return p, args.Error(1)
	}
	return args.Get(0).(Profile), nil
}

// mockCreateUseCase injects mock dependency into Controller layer.
type mockCreateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Create feature in the UseCase layer.
func (m *mockCreateUseCase) Execute(ID string, name string, lastname string, password string) error {
	args := m.Called(ID, name, lastname, password)
	return args.Error(0)
}

// mockUpdateUseCase injects mock dependency into Controller layer.
type mockUpdateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Update feature in the UseCase layer.
func (m *mockUpdateUseCase) Execute(p Profile) error {
	args := m.Called(p)
	return args.Error(0)
}