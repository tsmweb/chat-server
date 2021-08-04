package server

import (
	"context"
	"github.com/tsmweb/chat-service/server/message"
	"github.com/tsmweb/go-helper-api/kafka"
)

// HandleMessage
type HandleMessage interface {
	Execute(ctx context.Context, msg message.Message) *ErrorEvent
	Stop()
}

type handleMessage struct {
	encoder message.Encoder
	producer kafka.Producer
}

func NewHandleMessage(
	encoder message.Encoder,
	producer kafka.Producer,
) HandleMessage {
	return &handleMessage{
		encoder: encoder,
		producer: producer,
	}
}

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

func (h *handleMessage) Stop() {
	h.producer.Close()
}

// HandleGroupMessage
type HandleGroupMessage interface {
	Execute(msg message.Message, chMessage chan<- message.Message) *ErrorEvent
}

type handleGroupMessage struct {
	repository Repository
}

func NewHandleGroupMessage(repository Repository) HandleGroupMessage {
	return &handleGroupMessage{
		repository: repository,
	}
}

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
