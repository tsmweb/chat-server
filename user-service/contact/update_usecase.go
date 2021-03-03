package contact

import "time"

// UpdateUseCase updates a Contact, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(contact *Contact) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(r Repository) UpdateUseCase {
	return &updateUseCase{repository: r}
}

// Execute performs the update use case.
func (u *updateUseCase) Execute(contact *Contact) error {
	err := contact.Validate()
	if err != nil {
		return err
	}

	contact.UpdatedAt = time.Now()

	rows, err := u.repository.Update(contact)
	if err != nil {
		return err
	}
	if rows <= 0 {
		return ErrContactNotFound
	}

	return nil
}