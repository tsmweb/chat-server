package group

import (
	"context"
	"github.com/tsmweb/user-service/common"
	"time"
)

// UpdateUseCase updates a Group, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(ctx context.Context, group *Group) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(r Repository) UpdateUseCase {
	return &updateUseCase{repository: r}
}

// Execute performs the update use case.
func (u *updateUseCase) Execute(ctx context.Context, group *Group) error {
	err := group.Validate()
	if err != nil {
		return err
	}

	authID, err := u.checkPermission(ctx, group.ID)
	if err != nil {
		return err
	}

	group.UpdatedBy = authID
	group.UpdatedAt = time.Now().UTC()

	ok, err := u.repository.Update(ctx, group)
	if err != nil {
		return err
	}
	if !ok {
		return ErrGroupNotFound
	}

	// TODO notify members

	return nil
}

func (u *updateUseCase) checkPermission(ctx context.Context, groupID string) (string, error) {
	authID := ctx.Value(common.AuthContextKey).(string)

	ok, err := u.repository.IsGroupAdmin(ctx, groupID, authID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrOperationNotAllowed
	}

	return authID, nil
}
