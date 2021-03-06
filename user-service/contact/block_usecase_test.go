package contact

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestBlockUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.Background()

	t.Run("when use case fails with ErrUserNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewBlockUseCase(r)
		err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("when use case fails with ErrContactAlreadyBlocked", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Block", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(cerror.ErrRecordAlreadyRegistered).
			Once()

		uc := NewBlockUseCase(r)
		err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.Equal(t, ErrContactAlreadyBlocked, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewBlockUseCase(r)
		err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.NotNil(t, err)

		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Block", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		err = uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Block", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		uc := NewBlockUseCase(r)
		err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.Nil(t, err)
	})
}
