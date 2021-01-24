package contact

// GetPresenceUseCase returns a presence (online or offline) of the Contact by profileID and contactID,
// otherwise an error is returned.
type GetPresenceUseCase interface {
	Execute(profileID, contactID string) (PresenceType, error)
}

type getPresenceUseCase struct {
	repository Repository
}

// NewGetPresenceUseCase create a new instance of GetPresenceUseCase.
func NewGetPresenceUseCase(r Repository) GetPresenceUseCase {
	return &getPresenceUseCase{repository: r}
}

// Execute performs the use case to get presence.
func (u *getPresenceUseCase) Execute(profileID, contactID string) (PresenceType, error) {
	presence, err := u.repository.GetPresence(profileID, contactID)
	if err != nil {
		return presence, err
	}
	if presence == NotFound {
		return NotFound, ErrContactNotFound
	}

	return presence, nil
}
