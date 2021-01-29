package contact

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeleteUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrContactNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Delete", mock.Anything, mock.Anything).
			Return(0, nil).
			Once()
		uc := NewDeleteUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Equal(t, ErrContactNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Delete", mock.Anything, mock.Anything).
			Return(0, errors.New("error")).
			Once()
		uc := NewDeleteUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Delete", mock.Anything, mock.Anything).
			Return(1, nil).
			Once()
		uc := NewDeleteUseCase(r)
		err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Nil(t, err)
	})
}
