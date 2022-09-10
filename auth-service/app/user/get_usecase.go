package user

import (
	"context"
	"errors"

	"github.com/tsmweb/auth-service/common/service"
	"github.com/tsmweb/go-helper-api/cerror"
)

// GetUseCase returns a User by ID, otherwise an error is returned.
type GetUseCase interface {
	Execute(ctx context.Context, ID string) (*User, error)
}

type getUseCase struct {
	tag        string
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(repository Repository) GetUseCase {
	return &getUseCase{
		tag:        "user.GetUseCase",
		repository: repository,
	}
}

// Execute executes the get use case.
func (u *getUseCase) Execute(ctx context.Context, ID string) (*User, error) {
	user, err := u.repository.Get(ctx, ID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			service.Warn(ID, u.tag, err.Error())
			return nil, ErrUserNotFound
		}

		service.Error(ID, u.tag, err)
		return nil, err
	}
	if user == nil {
		service.Warn(ID, u.tag, err.Error())
		return nil, ErrUserNotFound
	}

	return user, nil
}
