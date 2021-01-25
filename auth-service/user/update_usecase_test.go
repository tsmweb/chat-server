package user

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
		user := &User{
			ID: "+5518999999999",
			Name: "",
			LastName: "Jobs",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(user)

		assert.Equal(t, ErrNameValidateModel, err)
	})

	t.Run("when use case fails with ErrUserNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(0, nil).
			Once()

		profile := &User{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(-1, errors.New("error")).
			Once()

		profile := &User{
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

		user := &User{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(user)

		assert.Nil(t, err)
	})
}
