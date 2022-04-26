package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/app/group"
)

// mockGroupGetUseCase injects mock dependency into handler layer.
type mockGroupGetUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupGetUseCase) Execute(ctx context.Context, groupID string) (*group.Group, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*group.Group), nil
}

// mockGroupGetAllUseCase injects mock dependency into handler layer.
type mockGroupGetAllUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupGetAllUseCase) Execute(ctx context.Context, userID string) ([]*group.Group, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*group.Group), nil
}

// mockGroupCreateUseCase injects mock dependency into handler layer.
type mockGroupCreateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupCreateUseCase) Execute(ctx context.Context, name, description, owner string) (string, error) {
	args := m.Called(ctx, name, description, owner)
	if args.Get(0) == "" {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

// mockGroupUpdateUseCase injects mock dependency into handler layer.
type mockGroupUpdateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupUpdateUseCase) Execute(ctx context.Context, group *group.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

// mockGroupDeleteUseCase injects mock dependency into handler layer.
type mockGroupDeleteUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupDeleteUseCase) Execute(ctx context.Context, groupID string) error {
	args := m.Called(ctx, groupID)
	return args.Error(0)
}

// mockGroupAddMemberUseCase injects mock dependency into handler layer.
type mockGroupAddMemberUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupAddMemberUseCase) Execute(ctx context.Context, groupID string, userID string, admin bool) error {
	args := m.Called(ctx, groupID, userID, admin)
	return args.Error(0)
}

// mockGroupRemoveMemberUseCase injects mock dependency into handler layer.
type mockGroupRemoveMemberUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupRemoveMemberUseCase) Execute(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

// mockGroupSetAdminUseCase injects mock dependency into handler layer.
type mockGroupSetAdminUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockGroupSetAdminUseCase) Execute(ctx context.Context, member *group.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}
