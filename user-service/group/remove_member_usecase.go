package group

import (
	"context"
	"fmt"
	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/user-service/common"
)

// RemoveMemberUseCase removes member from Group, otherwise an error is returned.
type RemoveMemberUseCase interface {
	Execute(ctx context.Context, groupID, userID string) error
}

type removeMemberUseCase struct {
	repository Repository
	encoder    EventEncoder
	producer   kafka.Producer
}

// NewRemoveMemberUseCase create a new instance of RemoveMemberUseCase.
func NewRemoveMemberUseCase(
	repository Repository,
	encoder EventEncoder,
	producer kafka.Producer,
) RemoveMemberUseCase {
	return &removeMemberUseCase{
		repository: repository,
		encoder:    encoder,
		producer:   producer,
	}
}

// Execute performs the remove member use case.
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

	if err = u.notifyMember(ctx, groupID, userID); err != nil {
		return &ErrEventNotification{Msg: err.Error()}
	}

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

func (u *removeMemberUseCase) notifyMember(ctx context.Context, groupID, userID string) error {
	event := NewEvent(groupID, userID, EventRemoveMember)
	epb, err := u.encoder.Marshal(event)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", groupID, userID)
	if err = u.producer.Publish(ctx, []byte(key), epb); err != nil {
		return err
	}
	return nil
}
