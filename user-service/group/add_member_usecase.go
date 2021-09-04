package group

import (
	"context"
	"errors"
	"fmt"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/user-service/common"
)

// AddMemberUseCase adds member to a Group, otherwise an error is returned.
type AddMemberUseCase interface {
	Execute(ctx context.Context, groupID string, userID string, admin bool) error
}

type addMemberUseCase struct {
	repository Repository
	encoder    EventEncoder
	producer   kafka.Producer
}

// NewAddMemberUseCase create a new instance of AddMemberUseCase.
func NewAddMemberUseCase(
	repository Repository,
	encoder EventEncoder,
	producer kafka.Producer,
) AddMemberUseCase {
	return &addMemberUseCase{
		repository: repository,
		encoder:    encoder,
		producer:   producer,
	}
}

// Execute performs the add member use case.
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

	if err = u.notifyMember(ctx, groupID, userID); err != nil {
		return &ErrEventNotification{Msg: err.Error()}
	}

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

func (u *addMemberUseCase) notifyMember(ctx context.Context, groupID, userID string) error {
	event := NewEvent(groupID, userID, EventAddMember)
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
