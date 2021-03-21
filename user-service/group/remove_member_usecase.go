package group

import (
	"context"
	"github.com/tsmweb/use-service/common"
)

// RemoveMemberUseCase removes member from Group, otherwise an error is returned.
type RemoveMemberUseCase interface {
	Execute(ctx context.Context, groupID, userID string) error
}

type removeMemberUseCase struct {
	repository Repository
}

// NewRemoveMemberUseCase create a new instance of RemoveMemberUseCase.
func NewRemoveMemberUseCase(r Repository) RemoveMemberUseCase {
	return &removeMemberUseCase{repository: r}
}

// Execute performs the creation use case.
func (u *removeMemberUseCase) Execute(ctx context.Context, groupID, userID string) error {
	err := u.checkPermission(ctx, groupID, userID)
	if err != nil {
		return err
	}

	ok, err := u.repository.RemoveMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrMemberNotFound
	}

	// TODO notify member

	return nil
}

func (u *removeMemberUseCase) checkPermission(ctx context.Context, groupID, userID string) error {
	authID := ctx.Value(common.AuthContextKey).(string)

	// the group owner cannot be removed.
	ok, err := u.repository.IsGroupOwner(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if ok {
		return ErrGroupOwnerCannotRemoved
	}

	// the member can leave the group
	if authID == userID {
		return nil
	}

	// the member can be deleted from the group by the administrator
	ok, err = u.repository.IsGroupAdmin(ctx, groupID, authID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrOperationNotAllowed
	}

	return nil
}