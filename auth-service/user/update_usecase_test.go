package user

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
		user := &User{
			ID: "+5518999999999",
			Name: "",
			LastName: "Jobs",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, user)

		assert.Equal(t, ErrNameValidateModel, err)
	})

	t.Run("when use case fails with ErrOperationNotAllowed", func(t *testing.T) {
		//t.Parallel()
		profile := &User{
			ID: "+5518977777777",
			Name: "Steve",
			LastName: "Jobs",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, profile)

		assert.Equal(t, ErrOperationNotAllowed, err)
	})

	t.Run("when use case fails with ErrUserNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		profile := &User{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, profile)

		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		profile := &User{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, profile)

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		user := &User{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, user)

		assert.Nil(t, err)
	})
}
