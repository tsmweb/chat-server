package profile

import (
	"errors"
	"github.com/tsmweb/helper-go/cerror"
)

type CreateUseCase interface {
	Execute(ID string, name string, lastname string, password string) error
}

type createUseCase struct {
	repository Repository
}

func NewCreateUseCase(repository Repository) CreateUseCase {
	return &createUseCase{repository}
}

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
