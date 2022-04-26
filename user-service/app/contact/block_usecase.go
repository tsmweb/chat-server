package contact

import (
	"context"
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/kafka"
	"time"
)

// BlockUseCase blocks a contact, otherwise an error is returned.
type BlockUseCase interface {
	Execute(ctx context.Context, userID, blockedUserID string) error
}

type blockUseCase struct {
	repository Repository
	encoder    EventEncoder
	producer   kafka.Producer
}

// NewBlockUseCase create a new instance of BlockUseCase.
func NewBlockUseCase(
	repository Repository,
	encoder EventEncoder,
	producer kafka.Producer,
) BlockUseCase {
	return &blockUseCase{
		repository: repository,
		encoder:    encoder,
		producer:   producer,
	}
}

// Execute perform the block use case.
func (u *blockUseCase) Execute(ctx context.Context, userID, blockedUserID string) error {
	ok, err := u.repository.ExistsUser(ctx, blockedUserID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	err = u.repository.Block(ctx, userID, blockedUserID, time.Now().UTC())
	if err != nil {
		if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
			return ErrContactAlreadyBlocked
		}
		return err
	}

	if err = u.notify(ctx, userID, blockedUserID); err != nil {
		return &ErrEventNotification{Msg: err.Error()}
	}

	return nil
}

func (u *blockUseCase) notify(ctx context.Context, userID, contactID string) error {
	event := NewEvent(userID, contactID, EventBlockUser)
	epb, err := u.encoder.Marshal(event)
	if err != nil {
		return err
	}

	if err = u.producer.Publish(ctx, []byte(userID), epb); err != nil {
		return err
	}
	return nil
}