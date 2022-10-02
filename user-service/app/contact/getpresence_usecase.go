package contact

import (
	"context"

	"github.com/tsmweb/user-service/common/service"
)

// GetPresenceUseCase returns a presence (online or offline) of the Contact by userID and contactID,
// otherwise an error is returned.
type GetPresenceUseCase interface {
	Execute(ctx context.Context, profileID, contactID string) (PresenceType, error)
}

type getPresenceUseCase struct {
	tag        string
	repository Repository
}

// NewGetPresenceUseCase create a new instance of GetPresenceUseCase.
func NewGetPresenceUseCase(r Repository) GetPresenceUseCase {
	return &getPresenceUseCase{
		tag:        "contact::GetPresenceUseCase",
		repository: r,
	}
}

// Execute performs the use case to get presence.
func (u *getPresenceUseCase) Execute(ctx context.Context, userID, contactID string) (PresenceType, error) {
	presence, err := u.repository.GetPresence(ctx, userID, contactID)
	if err != nil {
		service.Error(userID, u.tag, err)
		return presence, err
	}
	if presence == NotFound {
		return NotFound, ErrContactNotFound
	}

	return presence, nil
}
