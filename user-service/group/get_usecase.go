package group

import (
	"context"
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/use-service/common"
)

// GetUseCase returns a Group by groupID, otherwise an error is returned.
type GetUseCase interface {
	Execute(ctx context.Context, groupID string) (*Group, error)
}

type getUseCase struct {
	repository Repository
}

// NewGetUseCase create a new instance of GetUseCase.
func NewGetUseCase(r Repository) GetUseCase {
	return &getUseCase{repository: r}
}

// Execute performs the get use case.
func (u *getUseCase) Execute(ctx context.Context, groupID string) (*Group, error) {
	authID := ctx.Value(common.AuthContextKey).(string)

	group, err := u.repository.Get(ctx, groupID, authID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotFound
	}

	return group, nil
}