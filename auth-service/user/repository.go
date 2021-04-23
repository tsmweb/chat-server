package user

import "context"

// Reader interface
type Reader interface {
	Get(ctx context.Context, ID string) (*User, error)
}

// Writer user writer
type Writer interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) (bool, error)
}

// Repository interface for user data source.
type Repository interface {
	Reader
	Writer
}