package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/use-service/common"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		group := &Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "",
			Description: "Group of friends",
			Owner:       "+5518999999999",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, group)
		assert.Equal(t, ErrNameValidateModel, err)
	})

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		//t.Parallel()
		group := &Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Friends",
			Description: "Group of friends",
			Owner:       "+5518999999999",
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, group)
		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with ErrContactNotFound", func(t *testing.T) {
		//t.Parallel()
		group := &Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Friends",
			Description: "Group of friends",
			Owner:       "+5518999999999",
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Update", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, group)
		assert.Equal(t, ErrGroupNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		group := &Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Friends",
			Description: "Group of friends",
			Owner:       "+5518999999999",
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, group)
		assert.NotNil(t, err)

		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Update", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		uc = NewUpdateUseCase(r)
		err = uc.Execute(ctx, group)
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		group := &Group{
			ID:          "be49afd2ee890805c21ddd55879db1387aec9751",
			Name:        "Friends",
			Description: "Group of friends",
			Owner:       "+5518999999999",
		}

		r := new(mockRepository)
		r.On("IsGroupAdmin", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		r.On("Update", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, group)
		assert.Nil(t, err)
	})
}
