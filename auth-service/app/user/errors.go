package user

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrOperationNotAllowed = errors.New("operation not allowed")
)
