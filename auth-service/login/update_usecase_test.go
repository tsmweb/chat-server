package login

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/auth-service/profile"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestUpdateUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		l := Login{
			ID: "+5518999999999",
			Password: "",
		}

		r := new(mockRepository)
		uc := NewUpdateUseCase(r)
		err := uc.Execute(l)

		assert.Equal(t, ErrPasswordValidateModel, err)
	})

	t.Run("when use case fails with ErrProfileNotFound", func(t *testing.T) {
		//t.Parallel()
		l := Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(cerror.ErrNotFound).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(l)

		assert.Equal(t, profile.ErrProfileNotFound, err)
	})

	t.Run("when use case fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		l := Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(cerror.ErrInternalServer).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(l)

		assert.Equal(t, cerror.ErrInternalServer, err)
	})

	t.Run("when use case success", func(t *testing.T) {
		//t.Parallel()
		l := Login{
			ID: "+5518999999999",
			Password: "123456",
		}

		r := new(mockRepository)
		r.On("Update", mock.Anything).
			Return(nil).
			Once()

		uc := NewUpdateUseCase(r)
		err := uc.Execute(l)

		assert.Nil(t, err)
	})
}
