package profile

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// CreateUseCase creates a new Profile, otherwise an error is returned.
type CreateUseCase interface {
	Execute(ID string, name string, lastname string, password string) error
}

type createUseCase struct {
	repository Repository
}

// NewCreateUseCase create a new instance of CreateUseCase.
func NewCreateUseCase(repository Repository) CreateUseCase {
	return &createUseCase{repository}
}

// Execute executes the creation use case.
func (u *createUseCase) Execute(ID string, name string, lastname string, password string) error {
	p, err := NewProfile(ID, name, lastname, password)
	if err != nil {
		return err
	}

	err = u.repository.Create(p)
	if err != nil {
		if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
			return err
		} else {
			return cerror.ErrInternalServer
		}
	}

	return nil
}
