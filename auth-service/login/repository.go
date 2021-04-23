package login

import "context"

// Repository interface for login data source.
type Repository interface {
	Login(ctx context.Context, login *Login) (bool, error)
	Update(ctx context.Context, login *Login) (bool, error)
}
