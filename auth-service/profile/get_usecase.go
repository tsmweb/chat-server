package profile

import (
	"errors"
	"github.com/tsmweb/helper-go/cerror"
)

type GetUseCase interface {
	Execute(ID string) (Profile, error)
}

type getUseCase struct {
	repository Repository
}

func NewGetUseCase(repository Repository) GetUseCase {
	return &getUseCase{repository}
}

func (u *getUseCase) Execute(ID string) (Profile, error) {
	profile, err := u.repository.Get(ID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return profile, ErrProfileNotFound
		} else {
			return profile, cerror.ErrInternalServer
		}
	}

	return profile, nil
}
