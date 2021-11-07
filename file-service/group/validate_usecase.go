package group

import (
	"context"
)

// ValidateUseCase validates group file access.
type ValidateUseCase interface {
	Execute(ctx context.Context, groupID, userID string, isAdmin bool) error
}

type validateUseCase struct {
	repository Repository
}

// NewValidateUseCase create a new instance of ValidateUseCase.
func NewValidateUseCase(r Repository) ValidateUseCase {
	return &validateUseCase{repository: r}
}

// Execute performs group file access validation.
func (u *validateUseCase) Execute(ctx context.Context, groupID, userID string, isAdmin bool) error {
	ok, err := u.repository.ExistsGroup(ctx, groupID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrGroupNotFound
	}

	if isAdmin {
		ok, err = u.repository.IsGroupAdmin(ctx, groupID, userID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrOperationNotAllowed
		}
	} else {
		ok, err = u.repository.IsGroupMember(ctx, groupID, userID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrOperationNotAllowed
		}
	}

	return nil
}