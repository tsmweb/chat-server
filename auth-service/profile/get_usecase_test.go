package profile

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/helper-go/cerror"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestGetUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when repository fails", func(t *testing.T) {
		t.Parallel()

		r := new(mockRepository)
		defer r.AssertExpectations(t)

		r.On("Get", mock.Anything).
			Return(nil, errors.New("error")).
			Once()

		uc := NewGetUseCase(r)
		_, err := uc.Execute("+5518999999999")

		assert.NotNil(t, err)
	})

	t.Run("when repository fails with ErrProfileNotFound", func(t *testing.T) {
		t.Parallel()

		r := new(mockRepository)
		defer r.AssertExpectations(t)

		r.On("Get", mock.Anything).
			Return(nil, cerror.ErrNotFound).
			Once()

		uc := NewGetUseCase(r)
		_, err := uc.Execute("+5518999999999")

		assert.Equal(t, ErrProfileNotFound, err)
	})

	t.Run("when repository succeeds", func(t *testing.T) {
		t.Parallel()

		r := new(mockRepository)
		defer r.AssertExpectations(t)

		profile := Profile{
			ID:       "+5518999999999",
			Name:     "Steve",
			LastName: "Jobs",
		}

		r.On("Get", "+5518999999999").
			Return(profile, nil).
			Once()

		uc := NewGetUseCase(r)
		p, err := uc.Execute("+5518999999999")

		assert.Nil(t, err)
		assert.Equal(t, profile, p)
	})
}
