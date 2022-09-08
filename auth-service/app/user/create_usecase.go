package user

import (
	"context"
	"errors"

	"github.com/tsmweb/auth-service/common/service"
	"github.com/tsmweb/go-helper-api/cerror"
)

// CreateUseCase creates a new User, otherwise an error is returned.
type CreateUseCase interface {
	Execute(ctx context.Context, ID, name, lastname, password string) error
}

type createUseCase struct {
	tag        string
	repository Repository
}

// NewCreateUseCase create a new instance of CreateUseCase.
func NewCreateUseCase(repository Repository) CreateUseCase {
	return &createUseCase{
		tag:        "CreateUseCase",
		repository: repository,
	}
}

// Execute executes the creation use case.
func (u *createUseCase) Execute(ctx context.Context, ID, name, lastname, password string) error {
	user, err := NewUser(ID, name, lastname, password)
	if err != nil {
		return err
	}

	if err = u.repository.Create(ctx, user); err != nil {
		if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
			service.Warn(ID, u.tag, err.Error())
			return ErrUserAlreadyExists
		} else {
			service.Error(ID, u.tag, err)
			return err
		}
	}

	return nil
}
