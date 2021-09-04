package group

import (
	"context"
	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/user-service/common"
)

// DeleteUseCase delete a Group, otherwise an error is returned.
type DeleteUseCase interface {
	Execute(ctx context.Context, groupID string) error
}

type deleteUseCase struct {
	repository Repository
	encoder    EventEncoder
	producer   kafka.Producer
}

// NewDeleteUseCase create a new instance of DeleteUseCase.
func NewDeleteUseCase(
	repository Repository,
	encoder EventEncoder,
	producer kafka.Producer,
) DeleteUseCase {
	return &deleteUseCase{
		repository: repository,
		encoder:    encoder,
		producer:   producer,
	}
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

	if err = u.notifyMember(ctx, groupID); err != nil {
		return &ErrEventNotification{Msg: err.Error()}
	}

	return nil
}

func (u *deleteUseCase) notifyMember(ctx context.Context, groupID string) error {
	event := NewEvent(groupID, "", EventDeleteGroup)
	epb, err := u.encoder.Marshal(event)
	if err != nil {
		return err
	}

	if err = u.producer.Publish(ctx, []byte(groupID), epb); err != nil {
		return err
	}
	return nil
}
