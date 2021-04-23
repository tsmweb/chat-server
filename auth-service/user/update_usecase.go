package user

import (
	"context"
	"github.com/tsmweb/auth-service/common"
	"time"
)

// UpdateUseCase updates a User, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(ctx context.Context, profile *User) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(repository Repository) UpdateUseCase {
	return &updateUseCase{repository}
}

// Execute executes the update use case.
func (u *updateUseCase) Execute(ctx context.Context, user *User) error {
	err := user.Validate(UPDATE)
	if err != nil {
		return err
	}

	err = u.checkPermission(ctx, user.ID)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now().UTC()

	ok, err := u.repository.Update(ctx, user)
	if err != nil {
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
