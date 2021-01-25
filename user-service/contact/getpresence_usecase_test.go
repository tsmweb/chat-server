package contact

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetPresenceUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrContactNotFound", func(t *testing.T) {
		//t.Parallel()
		var presence PresenceType = NotFound

		r := new(mockRepository)
		r.On("GetPresence", mock.Anything, mock.Anything).
			Return(presence, nil).
			Once()

		uc := NewGetPresenceUseCase(r)
		_, err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Equal(t, ErrContactNotFound, err)
	})

	t.Run("when use case fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("GetPresence", mock.Anything, mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		uc := NewGetPresenceUseCase(r)
		_, err := uc.Execute("+5518999999999", "+5518977777777")

		assert.NotNil(t, err)
	})

	t.Run("when use case succeeds", func(t *testing.T) {
		//t.Parallel()
		var presence PresenceType = Online

		r := new(mockRepository)
		r.On("GetPresence", mock.Anything, mock.Anything).
			Return(presence, nil).
			Once()

		uc := NewGetPresenceUseCase(r)
		p, err := uc.Execute("+5518999999999", "+5518977777777")

		assert.Nil(t, err)
		assert.Equal(t, presence, p)
	})
}
