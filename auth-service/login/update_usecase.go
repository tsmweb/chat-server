package login

import (
	"github.com/tsmweb/auth-service/profile"
	"github.com/tsmweb/go-helper-api/cerror"
)

// UpdateUseCase updates the login password, otherwise an error will be returned.
type UpdateUseCase interface {
	Execute(login *Login) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(r Repository) UpdateUseCase {
	return &updateUseCase{
		repository: r,
	}
}

// Execute executes the update use case.
func (u *updateUseCase) Execute(login *Login) error {
	err := login.Validate()
	if err != nil {
		return err
	}

	err = login.ApplyHashPassword()
	if err != nil {
		return cerror.ErrInternalServer
	}

	rows, err := u.repository.Update(login)
	if err != nil {
		return err
	}
	if rows <= 0 {
		return profile.ErrProfileNotFound
	}

	return nil
}
