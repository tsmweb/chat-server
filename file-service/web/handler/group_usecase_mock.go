package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// mockGroupValidateUseCase injects mock dependency into UseCase layer.
type mockGroupValidateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupValidateUseCase) Execute(ctx context.Context, groupID, userID string, isAdmin bool) error {
	args := m.Called(ctx, groupID, userID, isAdmin)
	return args.Error(0)
}
