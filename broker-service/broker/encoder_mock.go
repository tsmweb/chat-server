package broker

import (
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
)

// mockMessageEncoder injects mock message.Encoder dependency.
type mockMessageEncoder struct {
	mock.Mock
}

// Marshal represents the simulated method for the Marshal feature in the message.Encoder layer.
func (m *mockMessageEncoder) Marshal(msg *message.Message) ([]byte, error) {
	args := m.Called(msg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), nil
}

// mockUserEncoder injects mock user.Encoder dependency.
type mockUserEncoder struct {
	mock.Mock
}

// Marshal represents the simulated method for the Marshal feature in the user.Encoder layer.
func (m *mockUserEncoder) Marshal(usr *user.User) ([]byte, error) {
	args := m.Called(usr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), nil
}