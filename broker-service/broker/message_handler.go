package broker

import (
	"context"
	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/go-helper-api/kafka"
)

// HandleMessage handles messages.
type HandleMessage interface {
	// Execute performs message handling.
	Execute(ctx context.Context, msg message.Message) *ErrorEvent
}

type handleMessage struct {
	userRepository user.Repository
	msgRepository  message.Repository
	queue          kafka.Kafka
}

// NewHandleMessage implements the HandleMessage interface.
func NewHandleMessage(
	userRepository user.Repository,
	msgRepository message.Repository,
	queue kafka.Kafka,
) HandleMessage {
	return &handleMessage{
		userRepository: userRepository,
		msgRepository:  msgRepository,
		queue:          queue,
	}
}

func (h *handleMessage) Execute(ctx context.Context, msg message.Message) *ErrorEvent {
	return nil
}
