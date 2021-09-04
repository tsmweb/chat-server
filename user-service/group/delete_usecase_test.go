package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/common"
	"testing"
)

func TestDeleteUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	encode := new(mockEventEncoder)
	encode.On("Marshal", mock.Anything).
		Return([]byte{}, nil)

	producer := new(mockProducer)
	producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewDeleteUseCase(r, encode, producer)
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

		uc := NewDeleteUseCase(r, encode, producer)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewDeleteUseCase(r, encode, producer)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.NotNil(t, err)

		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		r.On("Delete", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewDeleteUseCase(r, encode, producer)
		err = uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.NotNil(t, err)

		r.On("Delete", mock.Anything, mock.Anything).
			Return(true, nil)
		p := new(mockProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error"))

		uc = NewDeleteUseCase(r, encode, p)
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

		uc := NewDeleteUseCase(r, encode, producer)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751")
		assert.Nil(t, err)
	})
}
