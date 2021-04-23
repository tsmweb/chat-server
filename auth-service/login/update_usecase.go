package login

import (
	"context"
	"github.com/tsmweb/auth-service/common"
	"github.com/tsmweb/go-helper-api/cerror"
	"time"
)

// UpdateUseCase updates the login password, otherwise an error will be returned.
type UpdateUseCase interface {
	Execute(ctx context.Context, login *Login) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(r Repository) UpdateUseCase {
	return &updateUseCase{
		repository: r,
	}
}

// Execute executes the update use case.
func (u *updateUseCase) Execute(ctx context.Context, login *Login) error {
	err := login.Validate()
	if err != nil {
		return err
	}

	err = u.checkPermission(ctx, login.ID)
	if err != nil {
		return err
	}

	login.UpdatedAt = time.Now().UTC()

	err = login.ApplyHashPassword()
	if err != nil {
		return cerror.ErrInternalServer
	}

	ok, err := u.repository.Update(ctx, login)
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
