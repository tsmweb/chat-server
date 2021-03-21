package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/use-service/common"
	"testing"
)

func TestRemoveMemberUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	t.Run("when use case fails with ErrGroupOwnerCannotRemoved", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewRemoveMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777")
		assert.Equal(t, ErrGroupOwnerCannotRemoved, err)
	})

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewRemoveMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777")
		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with ErrMemberNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("RemoveMember", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewRemoveMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777")
		assert.Equal(t, ErrMemberNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewRemoveMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777")
		assert.NotNil(t, err)

		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewRemoveMemberUseCase(r)
		err = uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777")
		assert.NotNil(t, err)

		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("RemoveMember", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewRemoveMemberUseCase(r)
		err = uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777")
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds (member leaves the group)", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		r.On("RemoveMember", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewRemoveMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518999999999")
		assert.Nil(t, err)
	})

	t.Run("when use case succeeds (admin deletes group member)", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupOwner", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("RemoveMember", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewRemoveMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777")
		assert.Nil(t, err)
	})
}
