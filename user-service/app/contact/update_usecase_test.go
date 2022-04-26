package contact

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.Background()

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		contact := &Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
			UserID:   "",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, contact)

		assert.Equal(t, ErrUserIDValidateModel, err)
	})

	t.Run("when use case fails with ErrContactNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(false, nil).
			Once()

		contact := &Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
			UserID:   "+5518999999999",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, contact)

		assert.Equal(t, ErrContactNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(false, errors.New("error")).
			Once()

		contact := &Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
			UserID:   "+5518999999999",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, contact)

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything, mock.Anything).
			Return(true, nil).
			Once()

		contact := &Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
			UserID:   "+5518999999999",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(ctx, contact)

		assert.Nil(t, err)
	})

}
