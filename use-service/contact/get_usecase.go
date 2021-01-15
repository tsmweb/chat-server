package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// GetUseCase returns a Contact by profileID and contactID, otherwise an error is returned.
type GetUseCase interface {
	Execute(profileID, contactID string) (Contact, error)
}

type getUseCase struct {
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(r Repository) GetUseCase {
	return &getUseCase{repository: r}
}

// Execute executes the get use case.
func (u *getUseCase) Execute(profileID, contactID string) (Contact, error) {
	contact, err := u.repository.Get(profileID, contactID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return contact, ErrContactNotFound
		}
		return contact, cerror.ErrInternalServer
	}

	return contact, nil
}
