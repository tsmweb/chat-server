package contact

import (
	"context"
	"time"
)

// UpdateUseCase updates a Contact, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(ctx context.Context, contact *Contact) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(r Repository) UpdateUseCase {
	return &updateUseCase{repository: r}
}

// Execute performs the update use case.
func (u *updateUseCase) Execute(ctx context.Context, contact *Contact) error {
	err := contact.Validate()
	if err != nil {
		return err
	}

	contact.UpdatedAt = time.Now().UTC()

	ok, err := u.repository.Update(ctx, contact)
	if err != nil {
		return err
	}
	if !ok {
		return ErrContactNotFound
	}

	return nil
}
