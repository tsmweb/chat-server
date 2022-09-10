package contact

import (
	"context"
	"errors"

	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/user-service/common/service"
)

// GetAllUseCase returns a list of contacts by userID, otherwise an error is returned.
type GetAllUseCase interface {
	Execute(ctx context.Context, userID string) ([]*Contact, error)
}

type getAllUseCase struct {
	tag        string
	repository Repository
}

// NewGetAllUseCase create a new instance of GetAllUseCase.
func NewGetAllUseCase(r Repository) GetAllUseCase {
	return &getAllUseCase{
		tag:        "contact.GetAllUseCase",
		repository: r,
	}
}

// Execute performs the use case to get all.
func (u *getAllUseCase) Execute(ctx context.Context, userID string) ([]*Contact, error) {
	contacts, err := u.repository.GetAll(ctx, userID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return nil, ErrContactNotFound
		}
		service.Error(userID, u.tag, err)
		return nil, err
	}
	if len(contacts) == 0 {
		return nil, ErrContactNotFound
	}

	return contacts, nil
}
