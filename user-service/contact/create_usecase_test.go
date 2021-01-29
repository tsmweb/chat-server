package contact

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
		err := uc.Execute("+5518977777777", "Bill", "Gates", "")

		assert.Equal(t, ErrUserIDValidateModel, err)
	})

	t.Run("when use case fails with ErrUserNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsUser", mock.Anything).
			Return(false, nil).
			Once()

		uc := NewCreateUseCase(r)
		err := uc.Execute("+5518977777777", "Bill", "Gates", "+5518999999999")

		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("when use case fails with ErrContactAlreadyExists", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsUser", mock.Anything).
			Return(true, nil).
			Once()
		r.On("Create", mock.Anything).
			Return(cerror.ErrRecordAlreadyRegistered).
			Once()

		uc := NewCreateUseCase(r)
		err := uc.Execute("+5518977777777", "Bill", "Gates", "+5518999999999")

		assert.Equal(t, ErrContactAlreadyExists, err)
	})

	t.Run("when use case fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsUser", mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewCreateUseCase(r)
		err := uc.Execute("+5518977777777", "Bill", "Gates", "+5518999999999")

		assert.NotNil(t, err)

		r.On("ExistsUser", mock.Anything).
			Return(true, nil).
			Once()
		r.On("Create", mock.Anything).
			Return(errors.New("error")).
			Once()

		err = uc.Execute("+5518977777777", "Bill", "Gates", "+5518999999999")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsUser", mock.Anything).
			Return(true, nil).
			Once()
		r.On("Create", mock.Anything).
			Return(nil).
			Once()

		uc := NewCreateUseCase(r)
		err := uc.Execute("+5518977777777", "Bill", "Gates", "+5518999999999")

		assert.Nil(t, err)
	})

}
