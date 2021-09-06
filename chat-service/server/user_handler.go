package server

import (
	"context"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/chat-service/server/user"
	"github.com/tsmweb/go-helper-api/kafka"
)

// HandleUserStatus handles user status.
type HandleUserStatus interface {
	// Execute performs user status handling.
	Execute(ctx context.Context, userID string, status user.Status) *ErrorEvent

	// Close connections.
	Close()
}

type handleUserStatus struct {
	encoder  user.Encoder
	producer kafka.Producer
}

// NewHandleUserStatus implements the HandleUserStatus interface.
func NewHandleUserStatus(
	encoder user.Encoder,
	producer kafka.Producer,
) HandleUserStatus {
	return &handleUserStatus{
		encoder:  encoder,
		producer: producer,
	}
}

// Execute performs user status handling as: publish in topic kafka.
func (h *handleUserStatus) Execute(ctx context.Context, userID string, status user.Status) *ErrorEvent {
	serverID := "OFF"
	if status == user.Online {
		serverID = config.HostID()
	}

	u := user.NewUser(userID, status, serverID)
	upb, err := h.encoder.Marshal(u)
	if err != nil {
		return NewErrorEvent(userID, "HandleUserStatus.Execute()", err.Error())
	}

	if err = h.producer.Publish(ctx, []byte(userID), upb); err != nil {
		return NewErrorEvent(userID, "HandleUserStatus.Execute()", err.Error())
	}

	return nil
}

// Close connection with kafka producer.
func (h *handleUserStatus) Close() {
	h.producer.Close()
}
