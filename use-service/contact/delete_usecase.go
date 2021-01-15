package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// DeleteUseCase delete a Contact, otherwise an error is returned.
type DeleteUseCase interface {
	Execute(contact Contact) error
}

type deleteUseCase struct {
	repository Repository
}

// NewDeleteUseCase create a new instance of DeleteUseCase.
func NewDeleteUseCase(r Repository) DeleteUseCase {
	return &deleteUseCase{repository: r}
}

// Execute executes the creation use case.
func (u *deleteUseCase) Execute(contact Contact) error {
	err := u.repository.Delete(contact)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return ErrContactNotFound
		}
		return cerror.ErrInternalServer
	}

	return nil
}
