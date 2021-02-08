package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// GetAllUseCase returns a list of contacts by userID, otherwise an error is returned.
type GetAllUseCase interface {
	Execute(userID string) ([]*Contact, error)
}

type getAllUseCase struct {
	repository Repository
}

// NewGetAllUseCase create a new instance of GetAllUseCase.
func NewGetAllUseCase(r Repository) GetAllUseCase {
	return &getAllUseCase{repository: r}
}

// Execute performs the use case to get all.
func (u *getAllUseCase) Execute(userID string) ([]*Contact, error) {
	contacts, err := u.repository.GetAll(userID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return nil, ErrContactNotFound
		}
		return nil, err
	}
	if contacts == nil || len(contacts) == 0 {
		return nil, ErrContactNotFound
	}

	return contacts, nil
}
