package contact

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/common"
	"testing"
)

func TestUnblockUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.Background()

	encode := new(mockEventEncoder)
	encode.On("Marshal", mock.Anything).
		Return([]byte{}, nil)

	producer := new(common.MockKafkaProducer)
	producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	t.Run("when use case fails with ErrUserNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Unblock", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewUnblockUseCase(r, encode, producer)
		err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Unblock", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewUnblockUseCase(r, encode, producer)
		err := uc.Execute(ctx, "+5518999999999", "+5518977777777")
		assert.NotNil(t, err)

		r.On("Unblock", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		p := new(common.MockKafkaProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		uc = NewUnblockUseCase(r, encode, p)
		err = uc.Execute(ctx, "+5518999999999", "+5518977777777")
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Unblock", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewUnblockUseCase(r, encode, producer)
		err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.Nil(t, err)
	})
}
