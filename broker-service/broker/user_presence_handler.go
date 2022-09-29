package broker

import (
	"context"
	"fmt"

	"github.com/tsmweb/broker-service/broker/user"
)

// UserPresenceHandler handles user presence.
type UserPresenceHandler interface {
	// Execute performs user presence handling.
	Execute(ctx context.Context, usr user.User) error
}

type userPresenceHandler struct {
	userRepository user.Repository
}

// NewUserPresenceHandler implements the UserPresenceHandler interface.
func NewUserPresenceHandler(userRepository user.Repository) UserPresenceHandler {
	return &userPresenceHandler{
		userRepository: userRepository,
	}
}

func (h *userPresenceHandler) Execute(ctx context.Context, usr user.User) error {
	if err := h.userRepository.UpdateUserPresenceCache(ctx, usr.ID, usr.ServerID,
		usr.Status); err != nil {
		return fmt.Errorf("UserPresenceHandler.Execute(%s). Error: %v", usr.ID, err.Error())
	}
	return nil
}
