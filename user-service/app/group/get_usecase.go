package group

import (
	"context"
	"errors"

	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/user-service/common"
	"github.com/tsmweb/user-service/common/service"
)

// GetUseCase returns a Group by groupID, otherwise an error is returned.
type GetUseCase interface {
	Execute(ctx context.Context, groupID string) (*Group, error)
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
func (u *getUseCase) Execute(ctx context.Context, groupID string) (*Group, error) {
	authID := ctx.Value(common.AuthContextKey).(string)

	group, err := u.repository.Get(ctx, groupID, authID)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return nil, ErrGroupNotFound
		}
		service.Error(authID, u.tag, err)
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotFound
	}

	return group, nil
}
