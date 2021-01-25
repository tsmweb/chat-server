package user

import (
	"github.com/tsmweb/go-helper-api/util/hashutil"
)

// User data model
type User struct {
	ID       string
	Name     string
	LastName string
	Password string
}

// NewUser create a new User
func NewUser(ID, name, lastname, password string) (*User, error) {
	p := &User{
		ID:       ID,
		Name:     name,
		LastName: lastname,
		Password: password,
	}

	err := p.Validate(CREATE)
	if err != nil {
		return p, err
	}

	pwd, err := hashutil.HashSHA1(password)
	if err != nil {
		return p, err
	}
	p.Password = pwd

	return p, nil
}

// Validate model User.
func (p User) Validate(op Operation) error {
	if p.ID == "" {
		return ErrIDValidateModel
	}
	if p.Name == "" {
		return ErrNameValidateModel
	}
	if op == CREATE && p.Password == "" {
		return ErrPasswordValidateModel
	}

	return nil
}