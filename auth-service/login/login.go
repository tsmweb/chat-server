package login

import (
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"time"
)

// Login data model.
type Login struct {
	ID       string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewLogin create a new Login.
func NewLogin(ID, password string) (*Login, error) {
	l := &Login{
		ID: ID,
		Password: password,
		CreatedAt: time.Now().UTC(),
	}

	err := l.Validate()
	if err != nil {
		return l, err
	}

	if err = l.ApplyHashPassword(); err != nil {
		return l, err
	}

	return l, nil
}

// ApplyHashPassword hashes the password in plain text.
func (l *Login) ApplyHashPassword() error {
	pwd, err := hashutil.HashSHA1(l.Password)
	if err != nil {
		return err
	}

	l.Password = pwd
	return nil
}

// Validate model Login.
func (l *Login) Validate() error {
	if l.ID == "" {
		return ErrIDValidateModel
	}
	if l.Password == "" {
		return ErrPasswordValidateModel
	}

	return nil
}
