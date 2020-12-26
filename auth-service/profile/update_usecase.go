package profile

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

type UpdateUseCase interface {
	Execute(profile Profile) error
}

type updateUseCase struct {
	repository Repository
}

func NewUpdateUseCase(repository Repository) UpdateUseCase {
	return &updateUseCase{repository}
}

func (u *updateUseCase) Execute(profile Profile) error {
	err := profile.Validate(UPDATE)
	if err != nil {
		return err
	}

	err = u.repository.Update(profile)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return ErrProfileNotFound
		} else {
			return cerror.ErrInternalServer
		}
	}

	return nil
}
