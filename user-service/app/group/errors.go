package group

import (
	"errors"
)

var (
	ErrGroupNotFound           = errors.New("group not found")
	ErrUserNotFound            = errors.New("user not found")
	ErrMemberNotFound          = errors.New("member not found")
	ErrMemberAlreadyExists     = errors.New("member already exists")
	ErrOperationNotAllowed     = errors.New("operation not allowed")
	ErrGroupOwnerCannotRemoved = errors.New("group owner cannot be removed")
	ErrGroupOwnerCannotChanged = errors.New("group owner cannot be changed")
)

// ErrEventNotification error thrown if publishing an event results in an error.
type ErrEventNotification struct {
	Msg string
}

// Error implements interface Error
func (e *ErrEventNotification) Error() string {
	return e.Msg
}
