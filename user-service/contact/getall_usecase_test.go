package contact

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestGetAllUseCase_Execute(t *testing.T) {
	//t.Parallel()
	ctx := context.Background()

	t.Run("when use case fails with ErrContactNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("GetAll", mock.Anything, mock.Anything).
			Return(nil, cerror.ErrNotFound).
			Once()

		uc := NewGetAllUseCase(r)
		_, err := uc.Execute(ctx, "+5518999999999")

		assert.Equal(t, ErrContactNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("GetAll", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		uc := NewGetAllUseCase(r)
		_, err := uc.Execute(ctx, "+5518999999999")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		contacts := []*Contact{
			{
				ID:       "+5518977777777",
				Name:     "Bill",
				LastName: "Gates",
				UserID:   "+5518999999999",
			},
			{
				ID:       "+5518966666666",
				Name:     "Elon",
				LastName: "Musk",
				UserID:   "+5518999999999",
			},
		}

		r := new(mockRepository)
		r.On("GetAll", mock.Anything, mock.Anything).
			Return(contacts, nil).
			Once()

		uc := NewGetAllUseCase(r)
		cs, err := uc.Execute(ctx, "+5518999999999")

		assert.Nil(t, err)
		assert.Equal(t, contacts, cs)
	})
}
