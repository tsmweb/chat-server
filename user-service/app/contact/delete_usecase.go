package contact

import (
	"context"

	"github.com/tsmweb/user-service/common/service"
)

// DeleteUseCase delete a Contact, otherwise an error is returned.
type DeleteUseCase interface {
	Execute(ctx context.Context, userID, contactID string) error
}

type deleteUseCase struct {
	tag        string
	repository Repository
}

// NewDeleteUseCase create a new instance of DeleteUseCase.
func NewDeleteUseCase(r Repository) DeleteUseCase {
	return &deleteUseCase{
		tag:        "contact.DeleteUseCase",
		repository: r,
	}
}

// Execute performs the delete use case.
func (u *deleteUseCase) Execute(ctx context.Context, userID, contactID string) error {
	ok, err := u.repository.Delete(ctx, userID, contactID)
	if err != nil {
		service.Error(userID, u.tag, err)
		return err
	}
	if !ok {
		return ErrContactNotFound
	}

	return nil
}
