package login

import (
	"errors"
	"github.com/tsmweb/auth-service/helper/setting"
	"github.com/tsmweb/auth-service/profile"
	"github.com/tsmweb/go-helper-api/auth"
	"github.com/tsmweb/go-helper-api/cerror"
	"log"
)

// Service provides business rules for the controller layer.
type Service interface {
	Login(login Login) (*TokenAuth, error)
	Update(login Login) error
}

type service struct {
	repository Repository
	jwt        auth.JWT
}

// NewService creates a new instance of Service.
func NewService(repository Repository, jwt auth.JWT) Service {
	return &service{repository, jwt}
}

// Login returns a token if the credentials are valid, otherwise an error
// is returned.
func (s *service) Login(login Login) (*TokenAuth, error) {
	err := login.Validate()
	if err != nil {
		return nil, &cerror.ErrValidateModel{Msg: err.Error()}
	}

	ok, err := s.repository.Login(login)
	if err != nil {
		log.Println("[!] LoginService->Login " + err.Error())

		if errors.Is(err, cerror.ErrNotFound) {
			return nil, profile.ErrProfileNotFound
		} else {
			return nil, cerror.ErrInternalServer
		}
	}

	if !ok {
		return nil, cerror.ErrUnauthorized
	}

	token, err := s.jwt.GenerateToken(login.ID, setting.ExpireToken())
	if err != nil {
		log.Println("[!] LoginService->Login " + err.Error())
		return nil, cerror.ErrInternalServer
	}

	if len(token) == 0 {
		return nil, cerror.ErrUnauthorized
	}

	response := &TokenAuth{Token: token}
	return response, nil
}

// Update returns error equal to nil if the password update was successful,
// otherwise the generated error will be returned.
func (s *service) Update(login Login) error {
	err := login.Validate()
	if err != nil {
		return &cerror.ErrValidateModel{Msg: err.Error()}
	}

	err = s.repository.Update(login)
	if err != nil {
		log.Println("[!] LoginService->Update " + err.Error())

		if errors.Is(err, cerror.ErrNotFound) {
			return profile.ErrProfileNotFound
		} else {
			return cerror.ErrInternalServer
		}
	}

	return nil
}


