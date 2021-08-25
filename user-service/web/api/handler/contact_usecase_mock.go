package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/contact"
)

// mockContactGetUseCase injects mock dependency into Handler layer.
type mockContactGetUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockContactGetUseCase) Execute(ctx context.Context, userID, contactID string) (*contact.Contact, error) {
	args := m.Called(ctx, userID, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contact.Contact), nil
}

// mockContactGetAllUseCase injects mock dependency into Handler layer.
type mockContactGetAllUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockContactGetAllUseCase) Execute(ctx context.Context, userID string) ([]*contact.Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*contact.Contact), nil
}

// mockContactGetPresenceUseCase injects mock dependency into Handler layer.
type mockContactGetPresenceUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockContactGetPresenceUseCase) Execute(ctx context.Context, userID, contactID string) (contact.PresenceType, error) {
	args := m.Called(ctx, userID, contactID)
	if args.Error(1) != nil {
		return contact.NotFound, args.Error(1)
	}
	return args.Get(0).(contact.PresenceType), nil
}

// mockContactCreateUseCase injects mock dependency into Handler layer.
type mockContactCreateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockContactCreateUseCase) Execute(ctx context.Context, ID, name, lastname, profileID string) error {
	args := m.Called(ctx, ID, name, lastname, profileID)
	return args.Error(0)
}

// mockContactUpdateUseCase injects mock dependency into Handler layer.
type mockContactUpdateUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockContactUpdateUseCase) Execute(ctx context.Context, contact *contact.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

// mockContactDeleteUseCase injects mock dependency into Handler layer.
type mockContactDeleteUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockContactDeleteUseCase) Execute(ctx context.Context, userID, contactID string) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

// mockContactBlockUseCase injects mock dependency into Handler layer.
type mockContactBlockUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockContactBlockUseCase) Execute(ctx context.Context, userID, contactID string) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

// mockContactUnblockUseCase injects mock dependency into Handler layer.
type mockContactUnblockUseCase struct {
	mock.Mock
}

// Execute represents the simulated method for the Execute feature in the UseCase layer.
func (m *mockContactUnblockUseCase) Execute(ctx context.Context, userID, contactID string) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}
