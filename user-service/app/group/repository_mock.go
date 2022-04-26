package group

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// mockRepository injects mock dependency into UserCase layer.
type mockRepository struct {
	mock.Mock
}

// Get represents the simulated method for the Get feature in the Repository layer.
func (m *mockRepository) Get(ctx context.Context, groupID, userID string) (*Group, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Group), nil
}

// GetAll represents the simulated method for the GetAll feature in the Repository layer.
func (m *mockRepository) GetAll(ctx context.Context, userID string) ([]*Group, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Group), nil
}

// ExistsUser represents the simulated method for the ExistsUser feature in the Repository layer.
func (m *mockRepository) ExistsUser(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// ExistsGroup represents the simulated method for the ExistsGroup feature in the Repository layer.
func (m *mockRepository) ExistsGroup(ctx context.Context, groupID string) (bool, error) {
	args := m.Called(ctx, groupID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// IsGroupAdmin represents the simulated method for the IsGroupAdmin feature in the Repository layer.
func (m *mockRepository) IsGroupAdmin(ctx context.Context, groupID, userID string) (bool, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// IsGroupOwner represents the simulated method for the IsGroupOwner feature in the Repository layer.
func (m *mockRepository) IsGroupOwner(ctx context.Context, groupID, userID string) (bool, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// Create represents the simulated method for the Create feature in the Repository layer.
func (m *mockRepository) Create(ctx context.Context, g *Group) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}

// Update represents the simulated method for the Update feature in the Repository layer.
func (m *mockRepository) Update(ctx context.Context, g *Group) (bool, error) {
	args := m.Called(ctx, g)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// Delete represents the simulated method for the Delete feature in the Repository layer.
func (m *mockRepository) Delete(ctx context.Context, groupID string) (bool, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// AddMember represents the simulated method for the AddMember feature in the Repository layer.
func (m *mockRepository) AddMember(ctx context.Context, mb *Member) error {
	args := m.Called(ctx, mb)
	return args.Error(0)
}

// SetAdmin represents the simulated method for the SetMemberAdmin feature in the Repository layer.
func (m *mockRepository) SetAdmin(ctx context.Context, mb *Member) (bool, error) {
	args := m.Called(ctx, mb)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}

// RemoveMember represents the simulated method for the RemoveMember feature in the Repository layer.
func (m *mockRepository) RemoveMember(ctx context.Context, groupID, userID string) (bool, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}
