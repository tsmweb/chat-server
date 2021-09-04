package group

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required id"}
var ErrNameValidateModel = &cerror.ErrValidateModel{Msg: "required name"}
var ErrOwnerValidateModel = &cerror.ErrValidateModel{Msg: "required owner"}
var ErrGroupIDValidateModel = &cerror.ErrValidateModel{Msg: "required group_id"}
var ErrUserIDValidateModel = &cerror.ErrValidateModel{Msg: "required user_id"}
var ErrGroupNotFound = errors.New("group not found")
var ErrUserNotFound = errors.New("user not found")
var ErrMemberNotFound = errors.New("member not found")
var ErrMemberAlreadyExists = errors.New("member already exists")
var ErrOperationNotAllowed = errors.New("operation not allowed")
var ErrGroupOwnerCannotRemoved = errors.New("group owner cannot be removed")
var ErrGroupOwnerCannotChanged = errors.New("group owner cannot be changed")

// ErrEventNotification error thrown if publishing an event results in an error.
type ErrEventNotification struct {
	Msg string
}

// Error implements interface Error
func (e *ErrEventNotification) Error() string {
	return e.Msg
}