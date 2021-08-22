package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/auth-service/login"
)

// mockLoginUseCase injects mock dependency into Handler layer.
type mockLoginUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockLoginUseCase) Execute(ctx context.Context, ID, password string) (string, error) {
	args := m.Called(ctx, ID, password)
	if args.Get(1) != nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), nil
}

// mockLoginUpdateUseCase injects mock dependency into Handler layer.
type mockLoginUpdateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockLoginUpdateUseCase) Execute(ctx context.Context, l *login.Login) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}
