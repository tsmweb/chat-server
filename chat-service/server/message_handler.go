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
	encoder message.Encoder
	producer kafka.Producer
}

// NewHandleMessage implements the HandleMessage interface.
func NewHandleMessage(
	encoder message.Encoder,
	producer kafka.Producer,
) HandleMessage {
	return &handleMessage{
		encoder: encoder,
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

// Close connection with kafka producer.
func (h *handleMessage) Close() {
	h.producer.Close()
}

// HandleGroupMessage handles group messages.
type HandleGroupMessage interface {
	// Execute performs group message handling.
	Execute(msg message.Message, chMessage chan<- message.Message) *ErrorEvent
}

type handleGroupMessage struct {
	repository Repository
}

// NewHandleGroupMessage implements the HandleGroupMessage interface.
func NewHandleGroupMessage(repository Repository) HandleGroupMessage {
	return &handleGroupMessage{
		repository: repository,
	}
}

// Execute performs group message handling.
// Load group members into database and send in message channel.
func (h *handleGroupMessage) Execute(msg message.Message, chMessage chan<- message.Message) *ErrorEvent {
	users, err := h.repository.GetGroupMembers(msg.Group)
	if err != nil {
		return NewErrorEvent(msg.From, "HandleGroupMessage.Execute()", err.Error())
	}

	for _, user := range users {
		msg.To = user
		chMessage <- msg
	}

	return nil
}
