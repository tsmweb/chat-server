package login

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsmweb/auth-service/helper/common"
	"github.com/tsmweb/auth-service/helper/setting"
	"github.com/tsmweb/auth-service/profile"
	"github.com/tsmweb/go-helper-api/cerror"
	"testing"
)

func TestLoginUseCase_Execute(t *testing.T) {
	//t.Parallel()

	t.Run("when use case fails with ErrValidateModel", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		j := new(common.MockJWT)
		uc := NewLoginUseCase(r, j)
		_, err := uc.Execute("+5518999999999", "")

		assert.Equal(t, ErrPasswordValidateModel, err)
	})

	t.Run("when use case fails with ErrProfileNotFound", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Login", mock.Anything).
			Return(false, cerror.ErrNotFound).
			Once()
		j := new(common.MockJWT)

		uc := NewLoginUseCase(r, j)
		_, err := uc.Execute("+5518999999999", "123456")

		assert.Equal(t, profile.ErrProfileNotFound, err)
	})

	t.Run("when use case fails with ErrUnauthorized", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Login", mock.Anything).
			Return(false, nil).
			Once()
		j := new(common.MockJWT)

		uc := NewLoginUseCase(r, j)
		_, err := uc.Execute("+5518999999999", "123456")

		assert.Equal(t, cerror.ErrUnauthorized, err)

		r.On("Login", mock.Anything).
			Return(true, nil).
			Once()

		j.On("GenerateToken", "+5518999999999", setting.ExpireToken()).
			Return("", nil).
			Once()

		_, err = uc.Execute("+5518999999999", "123456")

		assert.Equal(t, cerror.ErrUnauthorized, err)
	})

	t.Run("when use case fails with ErrInternalServer", func(t *testing.T) {
		//t.Parallel()
		r := new(mockRepository)
		r.On("Login", mock.Anything).
			Return(false, errors.New("error")).
			Once()
		j := new(common.MockJWT)

		uc := NewLoginUseCase(r, j)
		_, err := uc.Execute("+5518999999999", "123456")

		assert.Equal(t, cerror.ErrInternalServer, err)

		r.On("Login", mock.Anything).
			Return(true, nil).
			Once()

		j.On("GenerateToken", "+5518999999999", setting.ExpireToken()).
			Return(nil, errors.New("error")).
			Once()

		_, err = uc.Execute("+5518999999999", "123456")

		assert.Equal(t, cerror.ErrInternalServer, err)
	})

	t.Run("when use case success", func(t *testing.T) {
		//t.Parallel()
		token := "A1B2C3D4E5F6"

		r := new(mockRepository)
		r.On("Login", mock.Anything).
			Return(true, nil).
			Once()

		j := new(common.MockJWT)
		j.On("GenerateToken", "+5518999999999", setting.ExpireToken()).
			Return(token, nil).
			Once()

		uc := NewLoginUseCase(r, j)
		tokenUC, err := uc.Execute("+5518999999999", "123456")

		assert.Nil(t, err)
		assert.Equal(t, token, tokenUC)
	})
}
