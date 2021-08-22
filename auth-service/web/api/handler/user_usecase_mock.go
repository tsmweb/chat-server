package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/auth-service/user"
)

// mockUserGetUseCase injects mock dependency into Handler layer.
type mockUserGetUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockUserGetUseCase) Execute(ctx context.Context, ID string) (*user.User, error) {
	args := m.Called(ctx, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), nil
}

// mockUserCreateUseCase injects mock dependency into Handler layer.
type mockUserCreateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockUserCreateUseCase) Execute(ctx context.Context, ID, name, lastname, password string) error {
	args := m.Called(ctx, ID, name, lastname, password)
	return args.Error(0)
}

// mockUserUpdateUseCase injects mock dependency into Handler layer.
type mockUserUpdateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockUserUpdateUseCase) Execute(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}
