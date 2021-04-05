package contact

import (
	"context"
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
	"time"
)

// BlockUseCase blocks a contact, otherwise an error is returned.
type BlockUseCase interface {
	Execute(ctx context.Context, userID, blockedUserID string) error
}

type blockUseCase struct {
	repository Repository
}

// NewBlockUseCase create a new instance of BlockUseCase.
func NewBlockUseCase(r Repository) BlockUseCase {
	return &blockUseCase{repository: r}
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

	return nil
}