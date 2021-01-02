package login

import (
	"errors"
	"github.com/tsmweb/auth-service/helper/setting"
	"github.com/tsmweb/auth-service/profile"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
)

// LoginUseCase returns a token if the credentials are valid, otherwise an error
// is returned.
type LoginUseCase interface {
	Execute(ID string, password string) (string, error)
}

type loginUseCase struct {
	repository Repository
	jwt        auth.JWT
}

// NewLoginUseCase create a new instance of LoginUseCase.
func NewLoginUseCase(repository Repository, jwt auth.JWT) LoginUseCase {
	return &loginUseCase{repository, jwt}
}

// Execute executes the login use case.
func (u *loginUseCase) Execute(ID string, password string) (string, error) {
	l, err := NewLogin(ID, password)
	if err != nil {
		return "", err
	}

	ok, err := u.repository.Login(l)
	if err != nil {
		if errors.Is(err, cerror.ErrNotFound) {
			return "", profile.ErrProfileNotFound
		} else {
			return "", cerror.ErrInternalServer
		}
	}

	if !ok {
		return "", cerror.ErrUnauthorized
	}

	token, err := u.jwt.GenerateToken(ID, setting.ExpireToken())
	if err != nil {
		return "", cerror.ErrInternalServer
	}

	if len(token) == 0 {
		return "", cerror.ErrUnauthorized
	}

	return token, nil
}
