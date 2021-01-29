package contact

// DeleteUseCase delete a Contact, otherwise an error is returned.
type DeleteUseCase interface {
	Execute(userID, contactID string) error
}

type deleteUseCase struct {
	repository Repository
}

// NewDeleteUseCase create a new instance of DeleteUseCase.
func NewDeleteUseCase(r Repository) DeleteUseCase {
	return &deleteUseCase{repository: r}
}

// Execute performs the delete use case.
func (u *deleteUseCase) Execute(userID, contactID string) error {
	rows, err := u.repository.Delete(userID, contactID)
	if err != nil {
		return err
	}
	if rows <= 0 {
		return ErrContactNotFound
	}

	return nil
}
