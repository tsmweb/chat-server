package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// GetAllUseCase returns a list of contacts by profileID, otherwise an error is returned.
type GetAllUseCase interface {
	Execute(profileID string) ([]Contact, error)
}

type getAllUseCase struct {
	repository Repository
}

// NewGetAllUseCase create a new instance of GetAllUseCase.
func NewGetAllUseCase(r Repository) GetAllUseCase {
	return &getAllUseCase{repository: r}
}

// Execute executes the get use case.
func (u *getAllUseCase) Execute(profileID string) ([]Contact, error) {
	contacts, err := u.repository.GetAll(profileID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return contacts, ErrContactNotFound
		}
		return contacts, cerror.ErrInternalServer
	}
	if len(contacts) == 0 {
		return contacts, ErrContactNotFound
	}

	return contacts, nil
}
