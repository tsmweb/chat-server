package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// CreateUseCase creates a new Contact, otherwise an error is returned.
type CreateUseCase interface {
	Execute(ID, name, lastname, profileID string) error
}

type createUseCase struct {
	repository Repository
}

// NewCreateUseCase create a new instance of CreateUseCase.
func NewCreateUseCase(r Repository) CreateUseCase {
	return &createUseCase{repository: r}
}

// Execute performs the creation use case.
func (u *createUseCase) Execute(ID, name, lastname, profileID string) error {
	c, err := NewContact(ID, name, lastname, profileID)
	if err != nil {
		return err
	}

	ok, err := u.repository.ExistsProfile(ID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrProfileNotFound
	}

	err = u.repository.Create(c)
	if err != nil {
		if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
			return ErrContactAlreadyExists
		}
		return err
	}

	return nil
}
