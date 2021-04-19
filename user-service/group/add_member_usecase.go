package group

import (
	"context"
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/user-service/common"
)

// AddMemberUseCase adds member to a Group, otherwise an error is returned.
type AddMemberUseCase interface {
	Execute(ctx context.Context, groupID string, userID string, admin bool) error
}

type addMemberUseCase struct {
	repository Repository
}

// NewAddMemberUseCase create a new instance of AddMemberUseCase.
func NewAddMemberUseCase(r Repository) AddMemberUseCase {
	return &addMemberUseCase{repository: r}
}

// Execute performs the creation use case.
func (u *addMemberUseCase) Execute(ctx context.Context, groupID string, userID string, admin bool) error {
	ok, err := u.repository.ExistsGroup(ctx, groupID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrGroupNotFound
	}

	ok, err = u.repository.ExistsUser(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	err = u.checkPermission(ctx, groupID)
	if err != nil {
		return err
	}

	member, err := NewMember(groupID, userID, admin)
	if err != nil {
		return err
	}

	err = u.repository.AddMember(ctx, member)
	if err != nil {
		if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
			return ErrMemberAlreadyExists
		}
		return err
	}

	// TODO notify member

	return nil
}

func (u *addMemberUseCase) checkPermission(ctx context.Context, groupID string) error {
	authID := ctx.Value(common.AuthContextKey).(string)

	ok, err := u.repository.IsGroupAdmin(ctx, groupID, authID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrOperationNotAllowed
	}

	return nil
}