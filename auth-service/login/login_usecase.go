package login

import (
	"context"
	"github.com/tsmweb/auth-service/config"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
)

// LoginUseCase returns a token if the credentials are valid, otherwise an error
// is returned.
type LoginUseCase interface {
	Execute(ctx context.Context, ID, password string) (string, error)
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
func (u *loginUseCase) Execute(ctx context.Context, ID, password string) (string, error) {
	l, err := NewLogin(ID, password)
	if err != nil {
		return "", err
	}

	ok, err := u.repository.Login(ctx, l)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", cerror.ErrUnauthorized
	}

	payload := map[string]interface{}{
		"id": ID,
	}

	token, err := u.jwt.GenerateToken(payload, config.ExpireToken())
	if err != nil {
		return "", err
	}

	if len(token) == 0 {
		return "", cerror.ErrUnauthorized
	}

	return token, nil
}
