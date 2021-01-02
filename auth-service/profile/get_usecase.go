package profile

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// GetUseCase returns a Profile by ID, otherwise an error is returned.
type GetUseCase interface {
	Execute(ID string) (Profile, error)
}

type getUseCase struct {
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(repository Repository) GetUseCase {
	return &getUseCase{repository}
}

// Execute executes the get use case.
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
