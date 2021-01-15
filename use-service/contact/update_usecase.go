package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// UpdateUseCase updates a Contact, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(contact Contact) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(r Repository) UpdateUseCase {
	return &updateUseCase{repository: r}
}

// Execute executes the update use case.
func (u *updateUseCase) Execute(contact Contact) error {
	err := contact.Validate()
	if err != nil {
		return err
	}

	err = u.repository.Update(contact)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return ErrContactNotFound
		}
		return cerror.ErrInternalServer
	}

	return nil
}