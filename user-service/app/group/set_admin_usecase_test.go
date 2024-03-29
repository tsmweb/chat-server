package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/common"
	"testing"
)

func TestSetAdminUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	encode := new(mockEventEncoder)
	encode.On("Marshal", mock.Anything).
		Return([]byte{}, nil)

	producer := new(common.MockKafkaProducer)
	producer.On("Publish", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		member := &Member{
			GroupID: "",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		r := new(mockRepository)
		uc := NewSetAdminUseCase(r, encode, producer)
		err := uc.Execute(ctx, member)
		assert.Equal(t, ErrGroupIDValidateModel, err)
	})

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		//t.Parallel()
		member := &Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewSetAdminUseCase(r, encode, producer)
		err := uc.Execute(ctx, member)
		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with ErrGroupOwnerCannotChanged", func(t *testing.T) {
		//t.Parallel()
		member := &Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewSetAdminUseCase(r, encode, producer)
		err := uc.Execute(ctx, member)
		assert.Equal(t, ErrGroupOwnerCannotChanged, err)
	})

	t.Run("when use case fails with ErrMemberNotFound", func(t *testing.T) {
		//t.Parallel()
		member := &Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		r.On("SetAdmin", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewSetAdminUseCase(r,encode, producer)
		err := uc.Execute(ctx, member)
		assert.Equal(t, ErrMemberNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		member := &Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewSetAdminUseCase(r, encode, producer)
		err := uc.Execute(ctx, member)
		assert.NotNil(t, err)

		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewSetAdminUseCase(r, encode, producer)
		err = uc.Execute(ctx, member)
		assert.NotNil(t, err)

		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil)
		r.On("SetAdmin", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewSetAdminUseCase(r, encode, producer)
		err = uc.Execute(ctx, member)
		assert.NotNil(t, err)

		r.On("SetAdmin", mock.Anything, mock.Anything).
			Return(true, nil)
		p := new(common.MockKafkaProducer)
		p.On("Publish", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("error"))

		uc = NewSetAdminUseCase(r, encode, p)
		err = uc.Execute(ctx, member)
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		member := &Member{
			GroupID: "be49afd2ee890805c21ddd55879db1387aec9751",
			UserID:  "+5518977777777",
			Admin:   true,
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		r.On("SetAdmin", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewSetAdminUseCase(r, encode, producer)
		err := uc.Execute(ctx, member)
		assert.Nil(t, err)
	})
}
