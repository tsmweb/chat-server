package profile

import (
	"github.com/tsmweb/helper-go/util/hashutil"
)

// Presenter data model
type Profile struct {
	ID       string
	Name     string
	LastName string
	Password string
}

// NewRouter create a new profile
func NewProfile(ID string, name string, lastname string, password string) (Profile, error) {
	p := Profile{
		ID:       ID,
		Name:     name,
		LastName: lastname,
		Password: password,
	}

	err := p.Validate(CREATE)
	if err != nil {
		return p, err
	}

	pwd, err := generatePassword(password)
	if err != nil {
		return p, err
	}
	p.Password = pwd

	return p, nil
}

// Validate model Presenter.
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

func generatePassword(raw string) (string ,error) {
	hash, err := hashutil.HashSHA1(raw)
	if err != nil {
		return "", err
	}

	return hash, nil
}
