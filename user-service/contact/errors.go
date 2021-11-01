package contact

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrUserNotFound = errors.New("user not found")
var ErrContactNotFound = errors.New("contact not found")
var ErrContactAlreadyBlocked = errors.New("contact already blocked")
var ErrContactAlreadyExists = errors.New("contact already exists")
var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required id"}
var ErrUserIDValidateModel = &cerror.ErrValidateModel{Msg: "required user_id"}

// ErrEventNotification error thrown if publishing an event results in an error.
type ErrEventNotification struct {
	Msg string
}

// Error implements interface Error
func (e *ErrEventNotification) Error() string {
	return e.Msg
}
