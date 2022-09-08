package user

import (
	"context"
	"time"

	"github.com/tsmweb/auth-service/common"
	"github.com/tsmweb/auth-service/common/service"
)

// UpdateUseCase updates a User, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(ctx context.Context, profile *User) error
}

type updateUseCase struct {
	tag        string
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(repository Repository) UpdateUseCase {
	return &updateUseCase{
		tag:        "UpdateUseCase",
		repository: repository,
	}
}

// Execute executes the update use case.
func (u *updateUseCase) Execute(ctx context.Context, user *User) error {
	err := user.Validate(UPDATE)
	if err != nil {
		return err
	}

	if err = u.checkPermission(ctx, user.ID); err != nil {
		service.Warn(user.ID, u.tag, err.Error())
		return err
	}

	user.UpdatedAt = time.Now().UTC()

	ok, err := u.repository.Update(ctx, user)
	if err != nil {
		service.Error(user.ID, u.tag, err)
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	return nil
}

func (u *updateUseCase) checkPermission(ctx context.Context, userID string) error {
	authID := ctx.Value(common.AuthContextKey).(string)
	// checks if the ID owns the data
	if userID != authID {
		return ErrOperationNotAllowed
	}

	return nil
}
