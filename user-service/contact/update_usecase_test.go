package contact

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	//t.Parallel()

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
		err := uc.Execute(contact)

		assert.Equal(t, ErrUserIDValidateModel, err)
	})

	t.Run("when use case fails with ErrContactNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(0, nil).
			Once()

		contact := &Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
			UserID:   "+5518999999999",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(contact)

		assert.Equal(t, ErrContactNotFound, err)
	})

	t.Run("when use case fails with Error", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(0, errors.New("error")).
			Once()

		contact := &Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
			UserID:   "+5518999999999",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(contact)

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(1, nil).
			Once()

		contact := &Contact{
			ID:       "+5518977777777",
			Name:     "Bill",
			LastName: "Gates",
			UserID:   "+5518999999999",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(contact)

		assert.Nil(t, err)
	})

}
