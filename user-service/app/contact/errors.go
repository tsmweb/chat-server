package contact

import (
	"errors"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrContactNotFound       = errors.New("contact not found")
	ErrContactAlreadyBlocked = errors.New("contact already blocked")
	ErrContactAlreadyExists  = errors.New("contact already exists")
)

// ErrEventNotification error thrown if publishing an event results in an error.
type ErrEventNotification struct {
	Msg string
}

// Error implements interface Error
func (e *ErrEventNotification) Error() string {
	return e.Msg
}
