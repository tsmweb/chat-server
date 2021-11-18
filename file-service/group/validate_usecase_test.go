package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestValidateUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("when use case fails with ErrGroupNotFound", func(t *testing.T) {
		r := new(mockRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewValidateUseCase(r)
		err := uc.Execute(ctx,
			"be49afd2ee890805c21ddd55879db1387aec9751",
			"+5518977777777",
			true)
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		r := new(mockRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewValidateUseCase(r)
		err := uc.Execute(ctx,
			"be49afd2ee890805c21ddd55879db1387aec9751",
			"+5518977777777",
			true)
		assert.Equal(t, ErrOperationNotAllowed, err)

		r.On("IsGroupMember", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		err = uc.Execute(ctx,
			"be49afd2ee890805c21ddd55879db1387aec9751",
			"+5518977777777",
			false)
		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with Error" , func(t *testing.T) {
		r := new(mockRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewValidateUseCase(r)
		err := uc.Execute(ctx,
			"be49afd2ee890805c21ddd55879db1387aec9751",
			"+5518977777777",
			true)
		assert.NotNil(t, err)

		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		err = uc.Execute(ctx,
			"be49afd2ee890805c21ddd55879db1387aec9751",
			"+5518977777777",
			true)
		assert.NotNil(t, err)

		r.On("IsGroupMember", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		err = uc.Execute(ctx,
			"be49afd2ee890805c21ddd55879db1387aec9751",
			"+5518977777777",
			false)
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		r := new(mockRepository)
		r.On("ExistsGroup", mock.Anything, mock.Anything).
			Return(true, nil)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewValidateUseCase(r)
		err := uc.Execute(ctx,
			"be49afd2ee890805c21ddd55879db1387aec9751",
			"+5518977777777",
			true)
		assert.Nil(t, err)

		r.On("IsGroupMember", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		err = uc.Execute(ctx,
			"be49afd2ee890805c21ddd55879db1387aec9751",
			"+5518977777777",
			false)
		assert.Nil(t, err)
	})

}