package user

import (
	"context"
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// CreateUseCase creates a new User, otherwise an error is returned.
type CreateUseCase interface {
	Execute(ctx context.Context, ID, name, lastname, password string) error
}

type createUseCase struct {
	repository Repository
}

// NewCreateUseCase create a new instance of CreateUseCase.
func NewCreateUseCase(repository Repository) CreateUseCase {
	return &createUseCase{repository}
}

// Execute executes the creation use case.
func (u *createUseCase) Execute(ctx context.Context, ID, name, lastname, password string) error {
	user, err := NewUser(ID, name, lastname, password)
	if err != nil {
		return err
	}

	err = u.repository.Create(ctx, user)
	if err != nil {
		if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
			return ErrUserAlreadyExists
		} else {
			return err
		}
	}

	return nil
}
