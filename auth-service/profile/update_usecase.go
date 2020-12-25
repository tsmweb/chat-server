package profile

import (
	"errors"
	"github.com/tsmweb/helper-go/cerror"
)

type UpdateUseCase struct {
	repository Repository
}

func NewUpdateUseCase(repository Repository) *UpdateUseCase {
	return &UpdateUseCase{repository}
}

func (u *UpdateUseCase) Execute(profile Profile) error {
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
