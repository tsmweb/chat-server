package login

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/auth-service/common"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518999999999",
			Password: "",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, l)

		assert.Equal(t, ErrPasswordValidateModel, err)
	})

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518977777777",
			Password: "123456",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, l)

		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with ErrUserNotFound", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, l)

		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, l)

		assert.NotNil(t, err)
	})

	t.Run("when use case success", func(t *testing.T) {
		//t.Parallel()
		l := &Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, l)

		assert.Nil(t, err)
	})
}
