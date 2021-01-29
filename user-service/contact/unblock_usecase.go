package contact

// UnblockUseCase unblocks a contact, otherwise an error is returned.
type UnblockUseCase interface {
	Execute(userID, contactID string) error
}

type unblockUseCase struct {
	repository Repository
}

// NewUnblockUseCase create a new instance of UnblockUseCase.
func NewUnblockUseCase(r Repository) UnblockUseCase {
	return &unblockUseCase{repository: r}
}

// Execute perform the unblock use case.
func (u *unblockUseCase) Execute(userID, contactID string) error {
	ok, err := u.repository.Unblock(userID, contactID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrContactNotFound
	}

	return nil
}