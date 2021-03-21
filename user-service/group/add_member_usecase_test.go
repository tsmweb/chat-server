package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/use-service/common"
	"testing"
)

func TestAddMemberUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewAddMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewAddMemberUseCase(r)
		err := uc.Execute(ctx, "", "+5518977777777", false)
		assert.Equal(t, ErrGroupIDValidateModel, err)
	})

	t.Run("when use case fails with ErrGroupNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewAddMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with ErrUserNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewAddMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("when use case fails with ErrMemberAlreadyExists", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("AddMember", mock.Anything, mock.Anything).
			Return(cerror.ErrRecordAlreadyRegistered).
			Once()

		uc := NewAddMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.Equal(t, ErrMemberAlreadyExists, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewAddMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.NotNil(t, err)

		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewAddMemberUseCase(r)
		err = uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.NotNil(t, err)

		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewAddMemberUseCase(r)
		err = uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.NotNil(t, err)

		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("AddMember", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		uc = NewAddMemberUseCase(r)
		err = uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("ExistsUser", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("AddMember", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		uc := NewAddMemberUseCase(r)
		err := uc.Execute(ctx, "be49afd2ee890805c21ddd55879db1387aec9751", "+5518977777777", false)
		assert.Nil(t, err)
	})
}
