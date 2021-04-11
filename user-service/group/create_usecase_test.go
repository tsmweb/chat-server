package group

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/user-service/common"
	"testing"
)

func TestCreateUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.WithValue(context.Background(), common.AuthContextKey, "+5518999999999")

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		uc := NewCreateUseCase(r)
		_, err := uc.Execute(ctx, "", "Group of friends", "+5518999999999")
		assert.Equal(t, ErrNameValidateModel, err)
	})

	t.Run("when use case fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Create", mock.Anything, mock.Anything).
			Return(errors.New("error")).
			Once()

		uc := NewCreateUseCase(r)
		_, err := uc.Execute(ctx, "Friends", "Group of friends", "+5518999999999")
		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Create", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		uc := NewCreateUseCase(r)
		ID, err := uc.Execute(ctx, "Friends", "Group of friends", "+5518999999999")
		assert.Nil(t, err)
		assert.NotNil(t, ID)
	})
}
