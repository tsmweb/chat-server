package broker

import (
	"context"
	"github.com/tsmweb/broker-service/broker/message"
)

// OfflineMessageHandler handles offline messages.
type OfflineMessageHandler interface {
	// Execute performs message handling.
	Execute(ctx context.Context, msg message.Message) *ErrorEvent
}

type offlineMessageHandler struct {
	msgRepository  message.Repository
}

// NewOfflineMessageHandler implements the OfflineMessageHandler interface.
func NewOfflineMessageHandler(msgRepository  message.Repository) OfflineMessageHandler {
	return &offlineMessageHandler{
		msgRepository: msgRepository,
	}
}

// Execute performs message handling.
func (h *offlineMessageHandler) Execute(ctx context.Context, msg message.Message) *ErrorEvent {
	if err := h.msgRepository.AddMessage(ctx, msg); err != nil {
		return NewErrorEvent(msg.From, "OfflineMessageHandler.AddMessage()", err.Error())
	}
	return nil
}
