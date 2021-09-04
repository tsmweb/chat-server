package group

import (
	"context"
	"fmt"
	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/user-service/common"
	"time"
)

// SetAdminUseCase sets a member as a group administrator, otherwise an error is returned.
type SetAdminUseCase interface {
	Execute(ctx context.Context, member *Member) error
}

type setAdminUseCase struct {
	repository Repository
	encoder    EventEncoder
	producer   kafka.Producer
}

// NewSetAdminUseCase create a new instance of SetAdminUseCase.
func NewSetAdminUseCase(
	repository Repository,
	encoder EventEncoder,
	producer kafka.Producer,
) SetAdminUseCase {
	return &setAdminUseCase{
		repository: repository,
		encoder:    encoder,
		producer:   producer,
	}
}

// Execute performs the set administrator use case.
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

	if err = u.notifyMember(ctx, member.GroupID, member.UserID, member.Admin); err != nil {
		return &ErrEventNotification{Msg: err.Error()}
	}

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

func (u *setAdminUseCase) notifyMember(ctx context.Context, groupID, userID string, isAdmin bool) error {
	var evt EventType = EventAddAdmin
	if isAdmin == false {
		evt = EventRemoveAdmin
	}

	event := NewEvent(groupID, userID, evt)
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
