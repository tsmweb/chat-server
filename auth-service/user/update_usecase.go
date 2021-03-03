package user

import "time"

// UpdateUseCase updates a User, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(profile *User) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(repository Repository) UpdateUseCase {
	return &updateUseCase{repository}
}

// Execute executes the update use case.
func (u *updateUseCase) Execute(user *User) error {
	err := user.Validate(UPDATE)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()

	rows, err := u.repository.Update(user)
	if err != nil {
		return err
	}
	if rows <= 0 {
		return ErrUserNotFound
	}

	return nil
}
