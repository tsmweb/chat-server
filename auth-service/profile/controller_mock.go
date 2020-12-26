package profile

import "github.com/stretchr/testify/mock"

// mockController injects mock dependency into Router layer.
type mockController struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Controller layer.
func (m *mockController) Get(ID string) (Presenter, error) {
	p := Presenter{}
	args := m.Called(ID)
	if args.Get(0) == nil {
		return p, args.Error(1)
	}
	return args.Get(0).(Presenter), nil
}

// Create represents the simulated method for the Create feature in the
// Controller layer.
func (m *mockController) Create(p Presenter) error {
	args := m.Called(p)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the
// Controller layer.
func (m *mockController) Update(p Presenter) error {
	args := m.Called(p)
	return args.Error(0)
}
