package group

import (
	"context"
	"time"

	"github.com/tsmweb/go-helper-api/kafka"
	"github.com/tsmweb/user-service/common"
	"github.com/tsmweb/user-service/common/service"
)

// UpdateUseCase updates a Group, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(ctx context.Context, group *Group) error
}

type updateUseCase struct {
	tag        string
	repository Repository
	encoder    EventEncoder
	producer   kafka.Producer
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(
	repository Repository,
	encoder EventEncoder,
	producer kafka.Producer,
) UpdateUseCase {
	return &updateUseCase{
		tag:        "group::UpdateUseCase",
		repository: repository,
		encoder:    encoder,
		producer:   producer,
	}
}

// Execute performs the update use case.
func (u *updateUseCase) Execute(ctx context.Context, group *Group) error {
	err := group.Validate(UPDATE)
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
		service.Error(authID, u.tag, err)
		return err
	}
	if !ok {
		return ErrGroupNotFound
	}

	if err = u.notifyMember(ctx, group.ID); err != nil {
		service.Error(authID, u.tag, err)
		return &ErrEventNotification{Msg: err.Error()}
	}

	return nil
}

func (u *updateUseCase) checkPermission(ctx context.Context, groupID string) (string, error) {
	authID := ctx.Value(common.AuthContextKey).(string)

	ok, err := u.repository.IsGroupAdmin(ctx, groupID, authID)
	if err != nil {
		service.Error(authID, u.tag, err)
		return "", err
	}
	if !ok {
		service.Warn(authID, u.tag, ErrOperationNotAllowed.Error())
		return "", ErrOperationNotAllowed
	}

	return authID, nil
}

func (u *updateUseCase) notifyMember(ctx context.Context, groupID string) error {
	event := NewEvent(groupID, "", EventUpdateGroup)
	epb, err := u.encoder.Marshal(event)
	if err != nil {
		return err
	}

	if err = u.producer.Publish(ctx, []byte(groupID), epb); err != nil {
		return err
	}
	return nil
}
