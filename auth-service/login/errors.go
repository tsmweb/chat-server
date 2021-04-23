package login

import (
	"errors"
	"github.com/tsmweb/go-helper-api/cerror"
)

var ErrIDValidateModel = &cerror.ErrValidateModel{Msg: "required id"}
var ErrPasswordValidateModel = &cerror.ErrValidateModel{Msg: "required password"}
var ErrUserNotFound = errors.New("user not found")
var ErrOperationNotAllowed = errors.New("operation not allowed")
