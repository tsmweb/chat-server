package contact

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUnblockUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrUserNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Unblock", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewUnblockUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Unblock", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewUnblockUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Unblock", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewUnblockUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Nil(t, err)
	})
}
