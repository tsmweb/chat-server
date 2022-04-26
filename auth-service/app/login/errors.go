package login

import (
	"errors"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrOperationNotAllowed = errors.New("operation not allowed")
)
