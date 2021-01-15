package contact

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestDeleteUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrContactNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Delete", mock.Anything).
			Return(cerror.ErrNotFound).
			Once()

		contact := Contact{
			ID: "+5518977777777",
			ProfileID: "+5518999999999",
		}

		uc := NewDeleteUseCase(r)
		err := uc.Execute(contact)

		assert.Equal(t, ErrContactNotFound, err)
	})

	t.Run("when use case fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Delete", mock.Anything).
			Return(errors.New("error")).
			Once()

		contact := Contact{
			ID: "+5518977777777",
			ProfileID: "+5518999999999",
		}

		uc := NewDeleteUseCase(r)
		err := uc.Execute(contact)

		assert.Equal(t, cerror.ErrInternalServer, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Delete", mock.Anything).
			Return(nil).
			Once()

		contact := Contact{
			ID: "+5518977777777",
			ProfileID: "+5518999999999",
		}

		uc := NewDeleteUseCase(r)
		err := uc.Execute(contact)

		assert.Nil(t, err)
	})
}
