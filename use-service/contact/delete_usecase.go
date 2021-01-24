package contact

// DeleteUseCase delete a Contact, otherwise an error is returned.
type DeleteUseCase interface {
	Execute(contact *Contact) error
}

type deleteUseCase struct {
	repository Repository
}

// NewDeleteUseCase create a new instance of DeleteUseCase.
func NewDeleteUseCase(r Repository) DeleteUseCase {
	return &deleteUseCase{repository: r}
}

// Execute performs the delete use case.
func (u *deleteUseCase) Execute(contact *Contact) error {
	rows, err := u.repository.Delete(contact)
	if err != nil {
		return err
	}
	if rows <= 0 {
		return ErrContactNotFound
	}

	return nil
}
