package group

import (
	"context"
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// GetAllUseCase returns a list of groups by userID, otherwise an error is returned.
type GetAllUseCase interface {
	Execute(ctx context.Context, userID string) ([]*Group, error)
}

type getAllUseCase struct {
	repository Repository
}

// NewGetAllUseCase create a new instance of GetAllUseCase.
func NewGetAllUseCase(r Repository) GetAllUseCase {
	return &getAllUseCase{repository: r}
}

// Execute performs the use case to get all.
func (u *getAllUseCase) Execute(ctx context.Context, userID string) ([]*Group, error) {
	groups, err := u.repository.GetAll(ctx, userID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	if groups == nil || len(groups) == 0 {
		return nil, ErrGroupNotFound
	}

	return groups, nil
}