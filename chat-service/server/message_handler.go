package server

import (
	"context"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/go-helper-api/kafka"
)

// HandleMessage handles messages.
type HandleMessage interface {
	// Execute performs message handling.
	Execute(ctx context.Context, msg message.Message) *ErrorEvent

	// Close connections.
	Close()
}

type handleMessage struct {
	encoder  message.Encoder
	producer kafka.Producer
}

// NewHandleMessage implements the HandleMessage interface.
func NewHandleMessage(
	encoder message.Encoder,
	producer kafka.Producer,
) HandleMessage {
	return &handleMessage{
		encoder:  encoder,
		producer: producer,
	}
}

// Execute performs message handling as: encode and publish in topic kafka.
func (h *handleMessage) Execute(ctx context.Context, msg message.Message) *ErrorEvent {
	mpb, err := h.encoder.Marshal(&msg)
	if err != nil {
		return NewErrorEvent(msg.From, "HandleMessage.Execute()", err.Error())
	}

	if err = h.producer.Publish(ctx, []byte(msg.ID), mpb); err != nil {
		return NewErrorEvent(msg.From, "HandleMessage.Execute()", err.Error())
	}

	return nil
}

// Close connection with kafka userProducer.
func (h *handleMessage) Close() {
	h.producer.Close()
}
