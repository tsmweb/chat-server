package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestCreateUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		uc := NewCreateUseCase(r)
		err := uc.Execute("+5518999999999", "Steve", "Jobs", "")

		assert.Equal(t, ErrPasswordValidateModel, err)
	})

	t.Run("when use case fails with ErrUserAlreadyExists", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Create", mock.Anything).
			Return(cerror.ErrRecordAlreadyRegistered).
			Once()

		uc := NewCreateUseCase(r)
		err := uc.Execute("+5518999999999", "Steve", "Jobs", "123456")

		assert.Equal(t, ErrUserAlreadyExists, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Create", mock.Anything).
			Return(errors.New("error")).
			Once()

		uc := NewCreateUseCase(r)
		err := uc.Execute("+5518999999999", "Steve", "Jobs", "123456")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Create", mock.Anything).
			Return(nil).
			Once()

		uc := NewCreateUseCase(r)
		err := uc.Execute("+5518999999999", "Steve", "Jobs", "123456")

		assert.Nil(t, err)
	})
}
