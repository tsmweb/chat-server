package group

import "github.com/stretchr/testify/mock"

// mockEventEncoder injects mock group.EventEncoder dependency.
type mockEventEncoder struct {
	mock.Mock
}

// Marshal represents the simulated method for the Marshal feature in the group.EventEncoder layer.
func (m *mockEventEncoder) Marshal(e *Event) ([]byte, error) {
	args := m.Called(e)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), nil
}
