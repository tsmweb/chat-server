package profile

// UpdateUseCase updates a Profile, otherwise an error is returned.
type UpdateUseCase interface {
	Execute(profile *Profile) error
}

type updateUseCase struct {
	repository Repository
}

// NewUpdateUseCase create a new instance of UpdateUseCase.
func NewUpdateUseCase(repository Repository) UpdateUseCase {
	return &updateUseCase{repository}
}

// Execute executes the update use case.
func (u *updateUseCase) Execute(profile *Profile) error {
	err := profile.Validate(UPDATE)
	if err != nil {
		return err
	}

	rows, err := u.repository.Update(profile)
	if err != nil {
		return err
	}
	if rows <= 0 {
		return ErrProfileNotFound
	}

	return nil
}
