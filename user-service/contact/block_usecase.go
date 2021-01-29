package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

// BlockUseCase blocks a contact, otherwise an error is returned.
type BlockUseCase interface {
	Execute(userID, contactID string) error
}

type blockUseCase struct {
	repository Repository
}

// NewBlockUseCase create a new instance of BlockUseCase.
func NewBlockUseCase(r Repository) BlockUseCase {
	return &blockUseCase{repository: r}
}

// Execute perform the block use case.
func (u *blockUseCase) Execute(userID, contactID string) error {
	ok, err := u.repository.ExistsUser(contactID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}

	_, err = u.repository.Block(userID, contactID)
	if err != nil {
		if errors.Is(err, cerror.ErrRecordAlreadyRegistered) {
			return ErrContactAlreadyBlocked
		}
		return err
	}

	return nil
}