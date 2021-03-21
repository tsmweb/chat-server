package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/use-service/common"
	"testing"
)

func TestDeleteUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewDeleteUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with ErrGroupNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Delete", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewDeleteUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewDeleteUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.NotNil(t, err)

		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Delete", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewDeleteUseCase(r)
		err = uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Delete", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewDeleteUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.Nil(t, err)
	})
}
