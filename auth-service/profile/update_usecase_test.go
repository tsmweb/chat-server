package profile

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when profile.Validate fails with ErrValidateModel", func(t *testing.T) {
		t.Parallel()

		profile := Profile{
			ID: "+5518999999999",
			Name: "",
			LastName: "Jobs",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.Equal(t, ErrNameValidateModel, err)
	})

	t.Run("when repository fails", func(t *testing.T) {
		t.Parallel()

		r := new(mockRepository)
		defer r.AssertExpectations(t)

		r.On("Update", mock.Anything).
			Return(errors.New("error")).
			Once()

		profile := Profile{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.NotNil(t, err)
	})

	t.Run("when repository fails with ErrNotFound", func(t *testing.T) {
		t.Parallel()

		r := new(mockRepository)
		defer r.AssertExpectations(t)

		r.On("Update", mock.Anything).
			Return(cerror.ErrNotFound).
			Once()

		profile := Profile{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.Equal(t, ErrProfileNotFound, err)
	})

	t.Run("when repository succeeds", func(t *testing.T) {
		t.Parallel()

		r := new(mockRepository)
		defer r.AssertExpectations(t)

		r.On("Update", mock.Anything).
			Return(nil).
			Once()

		profile := Profile{
			ID: "+5518999999999",
			Name: "Steve",
			LastName: "Jobs",
		}

		uc := NewUpdateUseCase(r)
		err := uc.Execute(profile)

		assert.Nil(t, err)
	})
}
