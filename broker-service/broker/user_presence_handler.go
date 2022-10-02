package broker

import (
	"context"

	"github.com/tsmweb/broker-service/broker/user"
	"github.com/tsmweb/broker-service/common/service"
)

// UserPresenceHandler handles user presence.
type UserPresenceHandler interface {
	// Execute performs user presence handling.
	Execute(ctx context.Context, usr user.User) error
}

type userPresenceHandler struct {
	tag            string
	userRepository user.Repository
}

// NewUserPresenceHandler implements the UserPresenceHandler interface.
func NewUserPresenceHandler(userRepository user.Repository) UserPresenceHandler {
	return &userPresenceHandler{
		tag:            "broker::UserPresenceHandler",
		userRepository: userRepository,
	}
}

func (h *userPresenceHandler) Execute(ctx context.Context, usr user.User) error {
	if err := h.userRepository.UpdateUserPresenceCache(ctx, usr.ID, usr.ServerID,
		usr.Status); err != nil {
		return service.FormatError(h.tag, err)
	}
	return nil
}
