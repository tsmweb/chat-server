package profile

import (
	"errors"
	"github.com/tsmweb/helper-go/cerror"
)

type GetUseCase struct {
	repository Repository
}

func NewGetUseCase(repository Repository) *GetUseCase {
	return &GetUseCase{repository}
}

func (u *GetUseCase) Execute(ID string) (Profile, error) {
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
