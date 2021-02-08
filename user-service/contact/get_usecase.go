package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// GetUseCase returns a Contact by userID and contactID, otherwise an error is returned.
type GetUseCase interface {
	Execute(userID, contactID string) (*Contact, error)
}

type getUseCase struct {
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(r Repository) GetUseCase {
	return &getUseCase{repository: r}
}

// Execute performs the get use case.
func (u *getUseCase) Execute(userID, contactID string) (*Contact, error) {
	contact, err := u.repository.Get(userID, contactID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return nil, ErrContactNotFound
		}
		return nil, err
	}
	if contact == nil {
		return nil, ErrContactNotFound
	}

	return contact, nil
}
