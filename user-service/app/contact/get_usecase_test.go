package contact

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestGetUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.Background()

	t.Run("when use case fails with ErrContactNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, cerror.ErrNotFound).
			Once()

		uc := NewGetUseCase(r)
		_, err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.Equal(t, ErrContactNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		uc := NewGetUseCase(r)
		_, err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		contact := &Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
			UserID:   "+5518999999999",
		}

		r := new(mockRepository)
		r.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(contact, nil).
			Once()

		uc := NewGetUseCase(r)
		c, err := uc.Execute(ctx, "+5518999999999", "+5518977777777")

		assert.Nil(t, err)
		assert.Equal(t, contact, c)
	})
}
