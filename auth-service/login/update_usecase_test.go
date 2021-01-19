package login

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/auth-service/profile"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518999999999",
			Password: "",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(l)

		assert.Equal(t, ErrPasswordValidateModel, err)
	})

	t.Run("when use case fails with ErrProfileNotFound", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(0, nil).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(l)

		assert.Equal(t, profile.ErrProfileNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(-1, errors.New("error")).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(l)

		assert.NotNil(t, err)
	})

	t.Run("when use case success", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(1, nil).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(l)

		assert.Nil(t, err)
	})
}
