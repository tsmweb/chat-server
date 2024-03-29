package contact

import (
	"context"

	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/user-service/common/service"
)

// UnblockUseCase unblocks a contact, otherwise an error is returned.
type UnblockUseCase interface {
	Execute(ctx context.Context, userID, blockedUserID string) error
}

type unblockUseCase struct {
	tag        string
	repository Repository
	encoder    EventEncoder
	producer   kafka.Producer
}

// NewUnblockUseCase create a new instance of UnblockUseCase.
func NewUnblockUseCase(
	repository Repository,
	encoder EventEncoder,
	producer kafka.Producer,
) UnblockUseCase {
	return &unblockUseCase{
		tag:        "contact::UnblockUseCase",
		repository: repository,
		encoder:    encoder,
		producer:   producer,
	}
}

// Execute perform the unblock use case.
func (u *unblockUseCase) Execute(ctx context.Context, userID, blockedUserID string) error {
	ok, err := u.repository.Unblock(ctx, userID, blockedUserID)
	if err != nil {
		service.Error(userID, u.tag, err)
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	if err = u.notify(ctx, userID, blockedUserID); err != nil {
		service.Error(userID, u.tag, err)
		return &ErrEventNotification{Msg: err.Error()}
	}

	return nil
}

func (u *unblockUseCase) notify(ctx context.Context, userID, contactID string) error {
	event := NewEvent(userID, contactID, EventUnblockUser)
	epb, err := u.encoder.Marshal(event)
	if err != nil {
		return err
	}

	if err = u.producer.Publish(ctx, []byte(userID), epb); err != nil {
		return err
	}
	return nil
}
