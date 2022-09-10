package login

import (
	"context"
	"time"

	"github.com/tsmweb/auth-service/common"
	"github.com/tsmweb/auth-service/common/service"
	"github.com/tsmweb/go-helper-api/cerror"
)

// UpdateUseCase updates the login password, otherwise an error will be returned.
type UpdateUseCase interface {
	Execute(ctx context.Context, login *Login) error
}

type updateUseCase struct {
	tag        string
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(r Repository) UpdateUseCase {
	return &updateUseCase{
		tag:        "login.UpdateUseCase",
		repository: r,
	}
}

// Execute executes the update use case.
func (u *updateUseCase) Execute(ctx context.Context, login *Login) error {
	err := login.Validate()
	if err != nil {
		return err
	}

	if err = u.checkPermission(ctx, login.ID); err != nil {
		service.Warn(login.ID, u.tag, err.Error())
		return err
	}

	login.UpdatedAt = time.Now().UTC()

	if err = login.ApplyHashPassword(); err != nil {
		service.Error(login.ID, u.tag, err)
		return cerror.ErrInternalServer
	}

	ok, err := u.repository.Update(ctx, login)
	if err != nil {
		service.Error(login.ID, u.tag, err)
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
