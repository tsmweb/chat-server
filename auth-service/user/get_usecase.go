package user

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// GetUseCase returns a User by ID, otherwise an error is returned.
type GetUseCase interface {
	Execute(ID string) (*User, error)
}

type getUseCase struct {
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(repository Repository) GetUseCase {
	return &getUseCase{repository}
}

// Execute executes the get use case.
func (u *getUseCase) Execute(ID string) (*User, error) {
	user, err := u.repository.Get(ID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
