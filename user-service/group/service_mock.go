package group

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// mockService injects mock dependency into Controller layer.
type mockService struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Service layer.
func (m *mockService) Get(ctx context.Context, groupID string) (*Group, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Group), nil
}

// GetAll represents the simulated method for the GetAll feature in the Service layer.
func (m *mockService) GetAll(ctx context.Context) ([]*Group, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Group), nil
}

// Create represents the simulated method for the Create feature in the Service layer.
func (m *mockService) Create(ctx context.Context, name, description, owner string) (string, error) {
	args := m.Called(ctx, name, description, owner)
	if args.Get(0) == "" {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

// Update represents the simulated method for the Update feature in the Service layer.
func (m *mockService) Update(ctx context.Context, group *Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

// Delete represents the simulated method for the Delete feature in the Service layer.
func (m *mockService) Delete(ctx context.Context, groupID string) error {
	args := m.Called(ctx, groupID)
	return args.Error(0)
}

// AddMember represents the simulated method for the AddMember feature in the Service layer.
func (m *mockService) AddMember(ctx context.Context, groupID string, userID string, admin bool) error {
	args := m.Called(ctx, groupID, userID, admin)
	return args.Error(0)
}

// RemoveMember represents the simulated method for the RemoveMember feature in the Service layer.
func (m *mockService) RemoveMember(ctx context.Context, groupID, userID string) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

// SetAdmin represents the simulated method for the SetAdmin feature in the Service layer.
func (m *mockService) SetAdmin(ctx context.Context, member *Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}
