package broker

import (
	"context"

	"github.com/tsmweb/broker-service/broker/message"
	"github.com/tsmweb/broker-service/common/service"
)

// OfflineMessageHandler handles offline messages.
type OfflineMessageHandler interface {
	// Execute performs message handling.
	Execute(ctx context.Context, msg message.Message) error
}

type offlineMessageHandler struct {
	tag           string
	msgRepository message.Repository
}

// NewOfflineMessageHandler implements the OfflineMessageHandler interface.
func NewOfflineMessageHandler(msgRepository message.Repository) OfflineMessageHandler {
	return &offlineMessageHandler{
		tag:           "broker::OfflineMessageHandler",
		msgRepository: msgRepository,
	}
}

// Execute performs message handling.
func (h *offlineMessageHandler) Execute(ctx context.Context, msg message.Message) error {
	if err := h.msgRepository.AddMessage(ctx, msg); err != nil {
		return service.FormatError(h.tag, err)
	}
	return nil
}
