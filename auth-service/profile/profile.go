package profile

import (
	"github.com/tsmweb/go-helper-api/util/hashutil"
)

// Profile data model
type Profile struct {
	ID       string
	Name     string
	LastName string
	Password string
}

// NewProfile create a new Profile
func NewProfile(ID, name, lastname, password string) (*Profile, error) {
	p := &Profile{
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

// Validate model Profile.
func (p Profile) Validate(op Operation) error {
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