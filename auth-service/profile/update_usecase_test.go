package profile

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		profile := &Profile{
			ID: "+5518999999999",
			Name: "",
			LastName: "Jobs",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.Equal(t, ErrNameValidateModel, err)
	})

	t.Run("when use case fails with ErrProfileNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(0, nil).
			Once()

		profile := &Profile{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.Equal(t, ErrProfileNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(-1, errors.New("error")).
			Once()

		profile := &Profile{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(1, nil).
			Once()

		profile := &Profile{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.Nil(t, err)
	})
}
