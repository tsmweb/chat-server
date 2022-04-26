package user

import (
	"context"
	"github.com/tsmweb/go-helper-api/cerror"
	"github.com/tsmweb/go-helper-api/util/hashutil"
	"time"
)

var (
	ErrIDValidateModel       = &cerror.ErrValidateModel{Msg: "required id"}
	ErrNameValidateModel     = &cerror.ErrValidateModel{Msg: "required name"}
	ErrPasswordValidateModel = &cerror.ErrValidateModel{Msg: "required password"}
)

// User data model
type User struct {
	ID        string
	Name      string
	LastName  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser create a new User
func NewUser(ID, name, lastname, password string) (*User, error) {
	p := &User{
		ID:        ID,
		Name:      name,
		LastName:  lastname,
		Password:  password,
		CreatedAt: time.Now().UTC(),
	}

	err := p.Validate(CREATE)
	if err != nil {
		return p, err
	}

	pwd, err := hashutil.HashSHA256(password)
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

// Repository interface for user data source.
type Repository interface {
	Get(ctx context.Context, ID string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) (bool, error)
}
