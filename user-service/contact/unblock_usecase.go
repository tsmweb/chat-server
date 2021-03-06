package contact

import "context"

// UnblockUseCase unblocks a contact, otherwise an error is returned.
type UnblockUseCase interface {
	Execute(ctx context.Context, userID, blockedUserID string) error
}

type unblockUseCase struct {
	repository Repository
}

// NewUnblockUseCase create a new instance of UnblockUseCase.
func NewUnblockUseCase(r Repository) UnblockUseCase {
	return &unblockUseCase{repository: r}
}

// Execute perform the unblock use case.
func (u *unblockUseCase) Execute(ctx context.Context, userID, blockedUserID string) error {
	ok, err := u.repository.Unblock(ctx, userID, blockedUserID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	return nil
}