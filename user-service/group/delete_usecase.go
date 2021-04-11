package group

import (
	"context"
	"github.com/tsmweb/user-service/common"
)

// DeleteUseCase delete a Group, otherwise an error is returned.
type DeleteUseCase interface {
	Execute(ctx context.Context, groupID string) error
}

type deleteUseCase struct {
	repository Repository
}

// NewDeleteUseCase create a new instance of DeleteUseCase.
func NewDeleteUseCase(r Repository) DeleteUseCase {
	return &deleteUseCase{repository: r}
}

// Execute performs the delete use case.
func (u *deleteUseCase) Execute(ctx context.Context, groupID string) error {
	authID := ctx.Value(common.AuthContextKey).(string)

	ok, err := u.repository.IsGroupAdmin(ctx, groupID, authID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrOperationNotAllowed
	}

	ok, err = u.repository.Delete(ctx, groupID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrGroupNotFound
	}

	// TODO notify member

	return nil
}