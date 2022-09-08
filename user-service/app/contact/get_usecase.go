package contact

import (
	"context"
	"errors"

	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/user-service/common/service"
)

// GetUseCase returns a Contact by userID and contactID, otherwise an error is returned.
type GetUseCase interface {
	Execute(ctx context.Context, userID, contactID string) (*Contact, error)
}

type getUseCase struct {
	tag        string
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(r Repository) GetUseCase {
	return &getUseCase{
		tag:        "GetUseCase",
		repository: r,
	}
}

// Execute performs the get use case.
func (u *getUseCase) Execute(ctx context.Context, userID, contactID string) (*Contact, error) {
	contact, err := u.repository.Get(ctx, userID, contactID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return nil, ErrContactNotFound
		}
		service.Error(userID, u.tag, err)
		return nil, err
	}
	if contact == nil {
		return nil, ErrContactNotFound
	}

	return contact, nil
}
