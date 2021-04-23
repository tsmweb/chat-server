package user

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// mockRepository injects mock dependency into UserCase layer.
type mockRepository struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Repository layer.
func (m *mockRepository) Get(ctx context.Context, ID string) (*User, error) {
	args := m.Called(ctx, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), nil
}

// Create represents the simulated method for the Create feature in the
// Repository layer.
func (m *mockRepository) Create(ctx context.Context, u *User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the
// Repository layer.
func (m *mockRepository) Update(ctx context.Context, u *User) (bool, error) {
	args := m.Called(ctx, u)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}