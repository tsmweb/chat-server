package contact

import "context"

// DeleteUseCase delete a Contact, otherwise an error is returned.
type DeleteUseCase interface {
	Execute(ctx context.Context, userID, contactID string) error
}

type deleteUseCase struct {
	repository Repository
}

// NewDeleteUseCase create a new instance of DeleteUseCase.
func NewDeleteUseCase(r Repository) DeleteUseCase {
	return &deleteUseCase{repository: r}
}

// Execute performs the delete use case.
func (u *deleteUseCase) Execute(ctx context.Context, userID, contactID string) error {
	ok, err := u.repository.Delete(ctx, userID, contactID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrContactNotFound
	}

	return nil
}
