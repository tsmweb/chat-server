package group

import (
	"context"
	"github.com/tsmweb/user-service/common"
	"time"
)

// SetAdminUseCase sets a member as a group administrator, otherwise an error is returned.
type SetAdminUseCase interface {
	Execute(ctx context.Context, member *Member) error
}

type setAdminUseCase struct {
	repository Repository
}

// NewSetAdminUseCase create a new instance of SetAdminUseCase.
func NewSetAdminUseCase(r Repository) SetAdminUseCase {
	return &setAdminUseCase{repository: r}
}

// Execute performs the creation use case.
func (u *setAdminUseCase) Execute(ctx context.Context, member *Member) error {
	err := member.Validate()
	if err != nil {
		return err
	}

	authID, err := u.checkPermission(ctx, member.GroupID, member.UserID)
	if err != nil {
		return err
	}

	member.UpdatedBy = authID
	member.UpdatedAt = time.Now().UTC()

	ok, err := u.repository.SetAdmin(ctx, member)
	if err != nil {
		return err
	}
	if !ok {
		return ErrMemberNotFound
	}

	// TODO notify member

	return nil
}

func (u *setAdminUseCase) checkPermission(ctx context.Context, groupID, userID string) (string, error) {
	authID := ctx.Value(common.AuthContextKey).(string)
	// the member cannot self-promote as an administrator.
	if authID == userID {
		return "", ErrOperationNotAllowed
	}

	// only group administrator can define a member to be admin.
	ok, err := u.repository.IsGroupAdmin(ctx, groupID, authID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrOperationNotAllowed
	}

	// the group owner cannot be changed.
	ok, err = u.repository.IsGroupOwner(ctx, groupID, userID)
	if err != nil {
		return "", err
	}
	if ok {
		return "", ErrGroupOwnerCannotChanged
	}

	return authID, nil
}