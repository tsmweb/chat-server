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
	Execute(ctx context.Context, userID string, status user.Status) error

	// Close connections.
	Close()
}

type handleUserStatus struct {
	encoder              user.Encoder
	userProducer         kafka.Producer
	userPresenceProducer kafka.Producer
}

// NewHandleUserStatus implements the HandleUserStatus interface.
func NewHandleUserStatus(
	encoder user.Encoder,
	userProducer kafka.Producer,
	userPresenceProducer kafka.Producer,
) HandleUserStatus {
	return &handleUserStatus{
		encoder:              encoder,
		userProducer:         userProducer,
		userPresenceProducer: userPresenceProducer,
	}
}

// Execute performs user status handling as: publish in topic kafka.
func (h *handleUserStatus) Execute(ctx context.Context, userID string, status user.Status) error {
	serverID := "OFF"
	if status == user.Online {
		serverID = config.HostID()
	}

	u := user.NewUser(userID, status, serverID)
	upb, err := h.encoder.Marshal(u)
	if err != nil {
		return err
	}

	if err = h.userProducer.Publish(ctx, []byte(userID), upb); err != nil {
		return err
	}

	if err = h.userPresenceProducer.Publish(ctx, []byte(userID), upb); err != nil {
		return err
	}

	return nil
}

// Close connection with kafka userProducer.
func (h *handleUserStatus) Close() {
	h.userProducer.Close()
	h.userPresenceProducer.Close()
}
