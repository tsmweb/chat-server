package contact

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestBlockUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrProfileNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsProfile", mock.Anything).
			Return(false, nil).
			Once()

		uc := NewBlockUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Equal(t, ErrProfileNotFound, err)
	})

	t.Run("when use case fails with ErrContactAlreadyBlocked", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsProfile", mock.Anything).
			Return(true, nil).
			Once()
		r.On("Block", mock.Anything, mock.Anything).
			Return(false, cerror.ErrRecordAlreadyRegistered).
			Once()

		uc := NewBlockUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Equal(t, ErrContactAlreadyBlocked, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsProfile", mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewBlockUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.NotNil(t, err)

		r.On("ExistsProfile", mock.Anything).
			Return(true, nil).
			Once()
		r.On("Block", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		err = uc.Execute("+5518999999999", "+5518977777777")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsProfile", mock.Anything).
			Return(true, nil).
			Once()
		r.On("Block", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewBlockUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Nil(t, err)
	})
}
