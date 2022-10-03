package broker

import (
	"context"

	"github.com/tsmweb/broker-service/broker/user"
)

// UserEventHandler handles user events.
type UserEventHandler interface {
	// Execute performs user event handling.
	Execute(ctx context.Context, evt user.Event) error
}

type userEventHandler struct {
	userRepository user.Repository
}

// NewUserEventHandler implements the UserEventHandler interface.
func NewUserEventHandler(userRepository user.Repository) UserEventHandler {
	return &userEventHandler{
		userRepository: userRepository,
	}
}

// Execute performs user event handling.
func (h *userEventHandler) Execute(ctx context.Context, evt user.Event) error {
	var isBlocked bool

	if evt.Event == user.EventBlockUser.String() {
		isBlocked = true
	} else if evt.Event == user.EventUnblockUser.String() {
		isBlocked = false
	} else {
		return nil
	}

	return h.userRepository.UpdateBlockedUserCache(
		ctx, evt.UserID, evt.ContactID, isBlocked)
}
