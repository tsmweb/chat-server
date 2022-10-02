package broker

import (
	"context"

	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/common/service"
)

// UserEventHandler handles user events.
type UserEventHandler interface {
	// Execute performs user event handling.
	Execute(ctx context.Context, evt user.Event) error
}

type userEventHandler struct {
	tag            string
	userRepository user.Repository
}

// NewUserEventHandler implements the UserEventHandler interface.
func NewUserEventHandler(userRepository user.Repository) UserEventHandler {
	return &userEventHandler{
		tag:            "broker::UserEventHandler",
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

	if err := h.userRepository.UpdateBlockedUserCache(
		ctx, evt.UserID, evt.ContactID, isBlocked,
	); err != nil {
		return service.FormatError(h.tag, err)
	}

	return nil
}
